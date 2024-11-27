
const body = document.querySelector('body') // тело страницы для добавления туда новых форм дедлайнов

const tasks= {
    item1: {
        width: '200px',
        height: '100px',
        backgroundColor: '#a8a8a8',
        text: 'item1'
    }
}

function createFormDeadline(taskName, deadline, priority) { // создание ui дедлайна
    const  div = document.createElement('div')
    div.style.cssText = `
      width: 400px;
      height: 100px;
      background-color: #353535;
      color: #FFFFFF
    `;
    div.innerHTML = `
        <h2>${taskName}</h2>
        <h3>${deadline}</h3>
        <h1>${priority}</h1>
    `;
    body.appendChild(div);
}

const buttonAddDeadline = document.getElementById('addDeadline');
buttonAddDeadline.addEventListener('click', function () // обработчик нажатия кнопки
{
    const taskName = document.getElementById('task-name').value;
    const description = document.getElementById('task-description').value;
    const deadline = document.getElementById('datetime-input').value;

    let priority = document.querySelector('input[name="slider"]:checked'); // поиск активного radio-элемента
    priority = priority ? parseInt(priority.id) : 0;

    createFormDeadline(taskName, deadline, priority); // создание визуального элемента дедлайна

    const formData = new FormData();
    formData.append("task-name", taskName);
    formData.append("task-description", description);
    formData.append("datetime-input", deadline);
    formData.append("priority", priority);

    fetch('/api/addDeadline', {method: 'POST', body: formData}) // асинхронная отправка данных на go-сервер
        .then(response => response.text()) // вывод ответа от сервера
        .catch(error => console.error('Ошибка:', error)); // вывод ошибки, в случае ошибки
});