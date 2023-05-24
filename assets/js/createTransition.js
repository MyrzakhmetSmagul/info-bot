let prevSelect;
let prevToState;
let prevFromState;
let labelHidden = true;
let language = document.getElementById("language");
let toState = document.getElementById("to-state");
let fromState = document.getElementById("from-state");
language.addEventListener("change", selectLanguage);

fromState.addEventListener("change", function () {
    let selectedOption = fromState.options[fromState.selectedIndex];

    let toStateOptions = toState.options;

    for (let i = 0; i < toStateOptions.length; i++) {
        if (selectedOption.value == toStateOptions[i].value) {
            toStateOptions[i].remove();
            break;
        }
    }
    if (prevFromState !== undefined) {
        toState.appendChild(prevFromState);
    }
    prevFromState = document.createElement("option");
    prevFromState = selectedOption;
    console.log(prevFromState);
})

toState.addEventListener("change", function () {
    let selectedOption = toState.options[toState.selectedIndex];

    let fromStateOptions = fromState.options;

    for (let i = 0; i < fromStateOptions.length; i++) {
        if (selectedOption.value == fromStateOptions[i].value) {
            fromStateOptions[i].remove();
            break;
        }
    }
    if (prevToState !== undefined) {
        fromState.appendChild(prevToState);
    }
    prevToState = document.createElement("option");
    prevToState = selectedOption;
    console.log(prevToState);
})

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