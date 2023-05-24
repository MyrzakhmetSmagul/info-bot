let prevSelect;
let prevSelect2;
let labelHidden = true;
let language = document.getElementById("language");
language.addEventListener("change", selectLanguage);

function selectLanguage() {
    if (labelHidden) {
        labelHidden = false;
        document.getElementById("message").classList.remove("hidden");
        document.getElementById("messageReply").classList.remove("hidden");
    }

    let prefix = language.value;
    let select = document.getElementById(prefix + "MsgGroup");
    let select2 = document.getElementById(prefix + "ReplyMsgGroup");

    select.classList.toggle("hidden");
    select2.classList.toggle("hidden");
    select.required = true;
    select2.required = true;
    if (prevSelect !== undefined) {
        prevSelect.classList.toggle("hidden");
        prevSelect2.classList.toggle("hidden");
        prevSelect.required = false;
        prevSelect2.required = false;
    }
    prevSelect = select;
    prevSelect2 = select2;
}