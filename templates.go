/*******************************************************************************
mailer: interfaces to send and fetch e-mail: https://github.com/mdhender/mailer

Copyright (c) 2021 Michael D Henderson.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
******************************************************************************/

package mailer

import (
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Template implements a template engine largely copied from
// https://gist.github.com/logrusorgru/abd846adb521a6fb39c7405f32fec0cf
type Template struct {
	*template.Template                  // root template
	dir                string           // root directory
	ext                string           // extension
	funcs              template.FuncMap // functions
}

// NewTemplate creates a new Template and performs the initial load of all template
// files that are in the template path.
// The dir argument is the directory to load templates from.
// The ext argument is extension of template files. It must include the
// leading dot. For example, ".gohtml," not "gohtml."
func NewTemplate(dir, ext string) (t *Template, err error) {
	if dir, err = filepath.Abs(dir); err != nil { // get absolute path
		return nil, err
	}
	t = &Template{dir: dir, ext: ext}
	if err = t.Load(); err != nil {
		return nil, err
	}
	return t, err
}

// Dir returns absolute path to directory with views
func (t *Template) Dir() string {
	return t.dir
}

// Ext returns extension of views
func (t *Template) Ext() string {
	return t.ext
}

// Funcs sets template functions
func (t *Template) Funcs(funcMap template.FuncMap) {
	t.Template = t.Template.Funcs(funcMap)
	t.funcs = funcMap
}

// Load all templates by walking the template path, finding all files
// that match the template extension, and loading each of them.
func (t *Template) Load() (err error) {
	// unnamed root template
	var root = template.New("")

	var walkFunc = func(path string, info os.FileInfo, err error) (_ error) {
		// handle walking error if any
		if err != nil {
			return err
		}

		// skip if it not a regular file.
		// yes, that means ignore symbolic links.
		if !info.Mode().IsRegular() {
			return
		}

		// filter by extension
		if filepath.Ext(path) != t.ext {
			return
		}

		// get path relative to the root of the templates
		var rel string
		if rel, err = filepath.Rel(t.dir, path); err != nil {
			return err
		}

		// name of a template is its relative path without extension
		rel = strings.TrimSuffix(rel, t.ext)

		// load the new template
		this := root.New(rel)
		if b, err := ioutil.ReadFile(path); err != nil {
			return err
		} else if _, err = this.Parse(string(b)); err != nil {
			return err
		}
		return nil
	}

	// walk the template tree, loading templates as we find them
	if err = filepath.Walk(t.dir, walkFunc); err != nil {
		return
	}

	// necessary for reloading
	if t.funcs != nil {
		root = root.Funcs(t.funcs)
	}

	t.Template = root // set or replace
	return
}

// Render the template using the data provided.
func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	return t.ExecuteTemplate(w, name, data)
}
