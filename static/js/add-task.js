const selector = e => document.querySelector(e) 

const taskInput = selector("#task-input") 
const submitButton = selector("#submit-btn")
const serverResponse = selector("#server-response")

const sendData = async (url, data) => {
    try {
        let useData = await fetch(url, {
            method: "POST",
            redirect: "follow",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });
        let parseRes = await useData.json();
        return parseRes;
    } catch (error) {
        return error;
    }
};


const grabInput = (e) => {
    try {
        if (e.target.value === "" || e.target.value.length === 0) {
            throw new Error("Target has nothing ")
        }
        taskInput.setAttribute("data_point", e.target.value)

    } catch (error) {
        console.log(error.message)
        return error
    }
}

const submitForm = e => {
    try {
        e.preventDefault() 
        console.log(taskInput.getAttribute("data_point"))
        sendData("/tasks", {
            title:taskInput.getAttribute("data_point")
        })
        .then((res) => {
            console.log(res);
            serverResponse.textContent = "";
            serverResponse.textContent = `Your title is ${res.title}`;
        })
        .catch((err) => (serverResponse.textContent = err.message));
    } catch (error) {
        console.log(error)
    }
}

taskInput.addEventListener("blur", grabInput)
submitButton.addEventListener("click", submitForm)
