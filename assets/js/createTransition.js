let fromState = document.getElementById("from-state");
let toState = document.getElementById("to-state");
let prevFromState;
let prevToState;

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
