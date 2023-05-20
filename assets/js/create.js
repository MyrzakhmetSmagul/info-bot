function createForm() {
    let pathname = window.location.pathname;
    console.log(pathname)
    let path = pathname.split("/")
    let formType = path[path.length - 1]
    switch (formType) {
        case "message-group":
            createMsgGroupForm()
            break;
        case "state":
            createStateForm()
    }
}

function createMsgGroupForm() {
// Create a form element
    let form = document.createElement("form");

// Create input elements for each field in the Message struct
    const lang = [["kz", "Қазақ"], ["ru", "Русский"], ["en", "English"]];
    for (let i = 0; i < 3; i++) {
        let prefix = lang[i][0];
        let triggerLabel = createLabel("Message trigger");
        let triggerField = createInputAbleElement("input", prefix + "MsgTrigger", "Message Trigger", "text");
        let textLabel = createLabel("Text");
        let textField = createInputAbleElement("textarea", prefix + "Text", "Text");
        let langField = createLabel("Language" + lang[i][i]);

        form.appendChild(triggerLabel);
        form.appendChild(triggerField);
        form.appendChild(textField);
        form.appendChild(textLabel);
        form.appendChild(langField);
        form.appendChild(document.createElement("br"));
    }

    let submit = document.createElement("input");
    submit.type = "submit";
    submit.textContent = "Создать группу сообщении"
    form.appendChild(submit)
    let formContainer = document.getElementById("form-container");
    formContainer.appendChild(form);
}

function createStateForm() {
    let form = document.createElement("form");

    let label = createLabel("Name");
    let field = createInputAbleElement("input", "State", "State", "text");
    let submit = document.createElement("input");
    submit.type = "submit";
    submit.textContent = "Создать группу сообщении"
    form.appendChild(label);
    form.appendChild(field);
    form.appendChild(submit);

    let formContainer = document.getElementById("form-container");
    formContainer.appendChild(form);
}

function createInputAbleElement(...args) {
    if (args.length < 3) {
        console.log("not enough arguments!")
        return
    }

    let element = document.createElement(args[0]);
    element.name = args[1];
    element.placeholder = args[2];
    if (args.length == 4) {
        element.type = args[3];
    }
    return element;
}

function createLabel(text) {
    let label = document.createElement("label");
    label.textContent = text;
    return label;
}