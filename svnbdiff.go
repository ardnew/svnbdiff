package main

import (
	"flag"
	"fmt"
	//"log"
	"os"
	//"os/exec"
	"strconv"
	"strings"
)

const (
	ErrorRevision = -1
	RevisionSep   = ","
)

type Revision int32
type RevisionPair struct {
	defined bool
	a       Revision
	b       Revision
}

func (r *Revision) String() string {
	return strconv.Itoa(int(*r))
}

func (r *Revision) Set(s string) error {
	*r = ErrorRevision
	v, err := strconv.ParseInt(s, 0, 32)
	if nil != err {
		return fmt.Errorf("revision number \"%s\" not an integer", s)
	}
	if v < 0 {
		return fmt.Errorf("revision number \"%d\" less than zero", v)
	}
	*r = Revision(v)
	return nil
}

func (p *RevisionPair) String() string {
	if (p.a != ErrorRevision) && (p.b != ErrorRevision) {
		return fmt.Sprintf("%ld%s%ld", p.a, RevisionSep, p.b)
	}
	if p.a != ErrorRevision {
		return fmt.Sprintf("%ld", p.a)
	}
	if p.b != ErrorRevision {
		return fmt.Sprintf("%ld", p.b)
	}
	return ""
}

func (p *RevisionPair) Set(s string) error {

	set := func(r *Revision, t []string, i uint8) error {
		if int(i) < len(t) {
			return r.Set(t[i])
		} else {
			*r = ErrorRevision
			return nil // not an error
		}
	}

	t := strings.Split(s, RevisionSep)
	if e := set(&p.a, t, 0); nil != e {
		return e
	}
	if e := set(&p.b, t, 1); nil != e {
		return e
	}

	p.defined = true
	return nil
}

type SVNPath struct {
	path string
	rev  Revision
	wc   bool
	head bool
}

func (p *SVNPath) String() string {
	a := []string{fmt.Sprintf("%d", p.rev)}
	if p.head {
		a = append(a, "HEAD")
	}
	if p.wc {
		a = append(a, "WC")
	}
	return fmt.Sprintf("%s@(%s)", p.path, strings.Join(a, ","))
}

var revision RevisionPair

func usage() {

	fmt.Println("")
	fmt.Println("-- USAGE -----------------------------------------------------------------------")
	fmt.Println("")
	fmt.Printf("    %s [-r REVS] PATH\n", os.Args[0])
	fmt.Println("")

	fmt.Println("")
	fmt.Println("-- OPTIONS ---------------------------------------------------------------------")
	fmt.Println("")
	flag.PrintDefaults()
	fmt.Println("")
	fmt.Println("\tif no options are provided:")
	fmt.Println("\t    1. the WC at PATH is compared to the HEAD revision")
	fmt.Println("\t        * PATH must be a valid SVN WC")
	fmt.Println("")
	fmt.Println("\tif a single revision REV is provided to option -r:")
	fmt.Println("\t    if PATH exists on local filesystem:")
	fmt.Println("\t        2. the revision at REV is compared to the WC at PATH")
	fmt.Println("\t            * PATH must be a valid SVN WC")
	fmt.Println("\t    if PATH does NOT exist on local filesystem:")
	fmt.Println("\t        3. the revision at REV is compared to the revision at HEAD")
	fmt.Println("\t            * PATH must be a fully-qualified SVN URL")
	fmt.Println("")
	fmt.Println("\tif two comma-separated revisions REV1,REV2 are provided to option -r:")
	fmt.Println("\t    4. the revision at REV1 is compared to the revision at REV2")
	fmt.Println("\t        * PATH must be a fully-qualified SVN URL")
	fmt.Println("")
}

func main() {

	flag.Usage = usage
	flag.Var(&revision, "r", "compare path at the specified SVN revision(s) `rev1[,rev2]`.")
	flag.Parse()

	if revision.defined {
		fmt.Printf("yes")
	} else {
		fmt.Printf("no")
	}

}
