package server

import (
	"html/template"
	"net/http"
)

func uploadModuleHandler(w http.ResponseWriter, r *http.Request) {
	data := createData("Upload Module", "box")
	templateFile := "module.html"

	data["url"] = "./joda"
	data["meta"] = template.HTML("<meta http-equiv=\"refresh\" content=\"5; url=/joda\" />")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `module`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, header, err := r.FormFile("module")
	if err != nil {
		addErrorMessage(data, "No module uploaded")
		executeTemplate(w, templateFile, data)
		return
	}
	defer file.Close()

	// Pass on file to JODA
	err = JodaInstance.UploadModule(file, header)
	if err != nil {
		addError(data, err)
		executeTemplate(w, templateFile, data)
		return
	}

	addSuccessMessage(data, "Module uploaded successfully")

	executeTemplate(w, templateFile, data)

}
