package addressbook

import (
	"fmt"
	"io"

	addressbook "github.com/MyWorkSpace/lets_Go/practice/protocols/addressBook"
)

func writePerson(w io.Writer, p *addressbook.Person) {
	fmt.Fprintln(w, "Person ID:", p.Id)
	fmt.Fprintln(w, "  Name:", p.Name)
	if p.Email != "" {
		fmt.Fprintln(w, "  E-mail address:", p.Email)
	}

	for _, pn := range p.Phones {
		switch pn.Type {
		case addressbook.Person_MOBILE:
			fmt.Fprint(w, "  Mobile phone #: ")
		case addressbook.Person_HOME:
			fmt.Fprint(w, "  Home phone #: ")
		case addressbook.Person_WORK:
			fmt.Fprint(w, "  Work phone #: ")
		}
		fmt.Fprintln(w, pn.Number)
	}
}

func listPeople(w io.Writer, book *addressbook.AddressBook) {
	for _, p := range book.People {
		writePerson(w, p)
	}
}

// Main reads the entire address book from a file and prints all the
// information inside.
