package main
import (
  "fmt"
  "regexp"
  "os"
  "path/filepath"
  "strings"
  "sort"
  "encoding/json"
)




type PG struct {
  FileName string   `json:"FileName"`
  Size     int64    `json:"Size"`
  Code   string     `json:"Code"`
}

func (p PG) String() string {
	return fmt.Sprintf("%s\t%d\t%s\n", p.Code, p.Size, p.FileName)
}

type BySize []PG;
func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Size < a[j].Size }

type BySizeDesc []PG;
func (a BySizeDesc) Len() int           { return len(a) }
func (a BySizeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySizeDesc) Less(i, j int) bool { return a[i].Size > a[j].Size }

var r *regexp.Regexp;
var pglist []PG;
var selectedPglist []PG;





func extractpgncode (fname string) string {
  m :=r.FindString(fname);
  return m;
}

func getsize(fname string) int64 {
  fi, e := os.Stat(fname);
  if e != nil {
    return -1
  }
  return fi.Size()
}

func isSelected(p PG) bool {
  for _,pg := range selectedPglist {
    if p.Code==pg.Code {
      return true;
    }
  }
  return false;
}

func main(){
  if (len(os.Args)<2) {
    fmt.Println("Lütfen klasör");
    return ;
  }
  r, _ = regexp.Compile("[^-]*$")

  fileList := []string{}
  filepath.Walk(os.Args[1], func(path string, f os.FileInfo, err error) error {
    if strings.Contains(path,"osd.")==true {
      fileList = append(fileList, path)
    }
    return nil
  })

  for _, file := range fileList {
    pglist=append(pglist,  PG {Code : extractpgncode(file) , FileName : file, Size : getsize(file)} );
  }

  sort.Sort(BySizeDesc(pglist));
  //fmt.Printf("\n%v\n", pglist);
  //fmt.Printf("%#v", pglist);
  for _, p := range pglist {
    if isSelected(p)==false {
      selectedPglist=append(selectedPglist,p);
    }
  }

  b, _ := json.MarshalIndent(selectedPglist, "", "  ")
  fmt.Println(string(b));

}
