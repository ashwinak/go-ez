package main

import (
        "bytes"
        "bufio"
        "io"
        "os"
        "log"
        "os/exec"
        "fmt"
        "syscall"
        "encoding/xml"
	"html/template"
	"net/http"
	"strings"
        "github.com/Juniper/go-ez/eznc"
        "github.com/Juniper/go-netconf/netconf"
)

type Todo struct {
	Title string
	Done  bool
        D string
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func handler(w http.ResponseWriter, r *http.Request) {

    tmpl := template.Must(template.ParseFiles("layout.html"))

    mkeys, ok := r.URL.Query()["mode"]
    if !ok || len(mkeys[0]) < 1 {
      log.Println("Url Param 'mode' is missing")
      return
    }
    mode := mkeys[0]
       
    fkeys, ok := r.URL.Query()["file_commands"]
    if !ok || len(fkeys[0]) < 1 {
      log.Println("Url Param 'file_commands' is missing")
      return
    }
    file_commands := fkeys[0]
    file, err := os.Open(file_commands)
      if err != nil {
        log.Fatal(err)
      }
    defer file.Close()
    scanner := bufio.NewScanner(file)
       
    if (mode == "Online") {
        jkeys, ok := r.URL.Query()["jhost"]
        if !ok || len(jkeys[0]) < 1 {
          log.Println("Url Param 'jhost' is missing")
          return
        }
        jhost := jkeys[0]
       
        ekeys, ok := r.URL.Query()["ehost"]
        if !ok || len(ekeys[0]) < 1 {
          log.Println("Url Param 'ehost' is missing")
          return
        }
        ehost := ekeys[0]

        ukeys, ok := r.URL.Query()["user"]
        if !ok || len(ukeys[0]) < 1 {
          log.Println("Url Param 'user' is missing")
          return
        }
        user := ukeys[0]
       
        pkeys, ok := r.URL.Query()["password"]
        if !ok || len(pkeys[0]) < 1 {
          log.Println("Url Param 'password' is missing")
          return
        }
        password := pkeys[0]
     
        es, err := eznc.Device(ehost, user, password)
        if err != nil {
          panic(err)
        }

        js, err := eznc.Device(jhost, user, password)
        if err != nil {
          panic(err)
        } 

      for scanner.Scan() {
        opcmd := scanner.Text()
        s := strings.Split(opcmd, "$")
        fmt.Println(s)
        junoscmd, evocmd := s[0], s[1]
        if s[1] == "" {
          evocmd = s[0]
        }
        EvoTags := RpcExecution(ehost, "evo", evocmd, es)
        JunosTags := RpcExecution(jhost, "junos", junoscmd, js)
	data := TodoPageData{
			PageTitle: "JUNOS/EVO XML tags differences",
			Todos: []Todo{
				{Title: "OpCommands to run on JUNOS/EVO(order)", Done: true, D: opcmd},
				{Title: "JunosHost", Done: true, D: jhost},
				{Title: "EvoHost", Done: true, D: ehost},
				{Title: "JUNOS-XML-TAGS", Done: true, D: JunosTags},
				{Title: "EVO-XML-TAGS", Done: true, D: EvoTags},
				{Title: "EVO-ADDITIVE-TAGS", Done: true, D: DiffXml("/var/tmp/evo.txt", "/var/tmp/junos.txt")}, 
				{Title: "EVO-MISSING-TAGS", Done: true, D: DiffXml("/var/tmp/junos.txt", "/var/tmp/evo.txt")}, 
			},
		}
	tmpl.Execute(w, data)
      }
    } else if (mode == "Offline") {
        for scanner.Scan() {
            opcmd := scanner.Text()
            s := strings.Split(opcmd, "$")
            junoscmd, junos_out_file, evo_out_file, evocmd := s[0], s[1], s[2], s[3]
            if s[3] == "" {
               evocmd = s[0]
            }
            EvoTags := FileContents(evo_out_file, "evo")
            JunosTags := FileContents(junos_out_file, "junos")
            data := TodoPageData{
                        PageTitle: "JUNOS/EVO XML tags differences",
                        Todos: []Todo{
                                {Title: "OpCommands to run on JUNOS/EVO(order)", Done: true, D: junoscmd + "$" + evocmd},
                                {Title: "JunosHost", Done: true, D: "Not Applicable"},
                                {Title: "EvoHost", Done: true, D: "Not Applicable"},
                                {Title: "JUNOS-XML-TAGS", Done: true, D: JunosTags},
                                {Title: "EVO-XML-TAGS", Done: true, D: EvoTags},
                                {Title: "EVO-ADDITIVE-TAGS", Done: true, D: DiffXml("/var/tmp/evo.txt", "/var/tmp/junos.txt")},
                                {Title: "EVO-MISSING-TAGS", Done: true, D: DiffXml("/var/tmp/junos.txt", "/var/tmp/evo.txt")},
                        },
            }
            tmpl.Execute(w, data)
        }
    }
}

func main() {
    /* Start http server */
    http.HandleFunc("/", handler)
    fmt.Println("Open URL: http://localhost:8080/")
    http.ListenAndServe(":8080", nil)
}

func Execute(output_buffer *bytes.Buffer, stack ...*exec.Cmd) (err error) {
	var error_buffer bytes.Buffer
	pipe_stack := make([]*io.PipeWriter, len(stack)-1)
	i := 0
	for ; i < len(stack)-1; i++ {
		stdin_pipe, stdout_pipe := io.Pipe()
		stack[i].Stdout = stdout_pipe
		stack[i].Stderr = &error_buffer
		stack[i+1].Stdin = stdin_pipe
		pipe_stack[i] = stdout_pipe
	}
	stack[i].Stdout = output_buffer
	stack[i].Stderr = &error_buffer

	if err := call(stack, pipe_stack); err != nil {
		log.Fatalln(string(error_buffer.Bytes()), err)
	}
	return err
}

func call(stack []*exec.Cmd, pipes []*io.PipeWriter) (err error) {
	if stack[0].Process == nil {
		if err = stack[0].Start(); err != nil {
			return err
		}
	}
	if len(stack) > 1 {
		if err = stack[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				pipes[0].Close()
				err = call(stack[1:], pipes[1:])
			}
		}()
	}
	return stack[0].Wait()
}

func DiffXml(host1 string, host2 string) (string) {
    if fileExists(host1) && fileExists(host2) {
        cmd := exec.Command("/usr/bin/grep", "-vf", host2, host1)
        out, _ := cmd.CombinedOutput()
        return string(out)
    } else {
        return "Invalid Data"
    }
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func FileContents(filename string, ost string) (string) {
    filerc, err := os.Open(filename)
    if err != nil{
      log.Fatal(err)
    }
    defer filerc.Close()
    buf := new(bytes.Buffer)
    buf.ReadFrom(filerc)
    contents := buf.String()
    ParseXML(contents, ost)
    file_name := "/var/tmp/" + ost + ".txt"
    cmd := exec.Command("/bin/cat", file_name)
    out, _ := cmd.CombinedOutput()

    return string(out)
}

func RpcExecution(host string, ost string, opcmd string, conn *netconf.Session) (string) {

    // eznc.Cli(<connection_handle>, "operational command", "format")
    reply, err := eznc.Cli(conn, opcmd, "xml")
    file_name := "/var/tmp/" + ost + ".txt"
    if err != nil {
        //panic(err)
        ret := "Unsupported command (or) Command execution fails"
        // Remove the intermediate files
        deleteFile(file_name, ost)
        return string(ret)
    }
    text_reply, err := eznc.Cli(conn, opcmd, "text")
    if err != nil {
        //panic(err)
        ret := "Unsupported command (or) Command execution fails"
        // Remove the intermediate files
        deleteFile(file_name, ost)
        return string(ret)
    }
 
    opcmd = strings.Replace(opcmd, " ", "_", -1)
    save_output_file := "/var/tmp/" + opcmd + "_" + ost
    /* Delete the file if it exists */
    if _, err := os.Stat(save_output_file); err == nil {
         os.Remove(save_output_file)
    }
    save_output_fileptr, err := os.OpenFile(save_output_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

    io.WriteString(save_output_fileptr, "\nXML OUTPUT: \n")
    io.WriteString(save_output_fileptr, reply.Data)

    io.WriteString(save_output_fileptr, "\nCLI OUTPUT: \n")
    io.WriteString(save_output_fileptr, text_reply.Data)
    
    save_output_fileptr.Close()

    ParseXML(reply.Data, ost)

    cmd := exec.Command("/bin/cat", file_name)
    out, _ := cmd.CombinedOutput()
    //defer conn.Close()
     
    return string(out)
    // fmt.Printf("Reply: %+v", reply.Data)
}

func deleteFile(path string, rel string) {
    // delete file
    os.Remove(path)
    fmt.Println("Invalid Command on " + rel)
}

func ParseXML(xmlStr string, os_type string) (err error) {

    r := bytes.NewBufferString(xmlStr)
    decoder := xml.NewDecoder(r)

    filename := "/var/tmp/" + os_type + ".txt"
    /* Delete the file if it exists */
    if _, err := os.Stat(filename); err == nil {
         os.Remove(filename)
    }

    /* Create new file and append to it repeatedly */
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
        if err != nil {
           panic(err)
    }

    /* Traverse the XML hierarchy */
    for {
	t, err := decoder.Token()
	if err == io.EOF {
	   // io.EOF is a successful end
	   break
	}
	if err != nil {
	   fmt.Printf("decoder.Token() failed with '%s'\n", err)
	   break
	}

	switch v := t.(type) {
	    case xml.StartElement:
            if _, err := os.Stat(filename); err == nil {
                /* file exists */
                tag := fmt.Sprintf("<%s>", v.Name.Local); 
                cmd := exec.Command("/usr/bin/grep", "-Fsq", tag, filename)
                var waitStatus syscall.WaitStatus
                cmd.CombinedOutput();
                waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
                if waitStatus.ExitStatus() != 0 {
                   fmt.Fprintf(f, "<%s>\r\n", v.Name.Local)
                   defer f.Close()
                }
            }  else {   
                fmt.Fprintf(f, "<%s>\r\n", v.Name.Local)
                defer f.Close()
            }
        }
    }
    return
}
