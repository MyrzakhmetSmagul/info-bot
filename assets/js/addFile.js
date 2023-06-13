let selected = false;
let fileType = document.getElementById("fileType");
let fileUpload = document.getElementById("fileUpload");
let fileName = document.getElementById("fileName");
let deleteFileButton = document.getElementById("deleteFile");
let submitButton = document.getElementById("submit");
fileType.addEventListener("change", changeType);
fileUpload.addEventListener("change", validateFile);
deleteFileButton.addEventListener("click", removeFile);

function changeType() {
    if (!selected) {
        selected = true;
        document.getElementById("fileLabel").classList.toggle("hidden");
        document.getElementById("fileUpload").classList.toggle("hidden");
    }

    let selectedOption = fileType.options[fileType.selectedIndex];
    switch (selectedOption.value) {
        case "pdf":
            fileUpload.setAttribute("accept", ".pdf");
            break;
        case "video":
            fileUpload.setAttribute("accept", "video/*")
            break;
        case "photo":
            fileUpload.setAttribute("accept", "image/*")
            break;
    }
    document.getElementById("fileTypeValue").value = selectedOption.value
}

function validateFile() {
    let file = fileUpload.files[0];
    let selectedOption = fileType.options[fileType.selectedIndex];
    switch (selectedOption.value) {
        case "video":
            if (file.size > 30 * 1024 * 1024) {
                alert("File size exceeds the limit of 30MB");
                fileUpload.value = ""; // Clear the selected file
                return;
            }
            break;
        default:
            if (file.size > 5 * 1024 * 1024) {
                alert("File size exceeds the limit of 5MB");
                fileUpload.value = ""; // Clear the selected file
                return;
            }
    }
    if (fileUpload.files.length > 0) {
        fileType.disabled = true;
        console.log(document.getElementById("fileTypeValue").value)
        fileUpload.classList.toggle("hidden");
        fileName.classList.toggle("hidden");
        fileName.textContent = `Имя файла: ${file.name}`;
        deleteFileButton.classList.toggle("hidden");
        submitButton.classList.toggle("hidden");
    }
}

function removeFile() {
    fileUpload.value = ""
    fileUpload.classList.toggle("hidden");
    fileName.classList.toggle("hidden");
    deleteFileButton.classList.toggle("hidden");
    submitButton.classList.toggle("hidden");
    fileType.disabled = false;
    console.log(fileType.disabled)

}

let prevSelect;
let labelHidden = true;
let language = document.getElementById("language");
language.addEventListener("change", selectLanguageAddFile);

function selectLanguageAddFile() {
    if (labelHidden) {
        labelHidden = false;
        document.getElementById("message").classList.remove("hidden");
    }

    let prefix = language.value;
    let select = document.getElementById(prefix + "MsgGroup");

    select.classList.toggle("hidden");
    select.required = true;
    if (prevSelect !== undefined) {
        prevSelect.classList.toggle("hidden");
        prevSelect.required = false;
    }
    prevSelect = select;
}