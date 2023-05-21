let prevSelect;
let labelHidden = true;
let language = document.getElementById("language");
language.addEventListener("change", selectLanguage);

function selectLanguage() {
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