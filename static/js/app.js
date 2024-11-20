
const body = document.querySelector('body') // тело страницы для добавления туда новых форм дедлайнов

const tasks= {
    item1: {
        width: '200px',
        height: '100px',
        backgroundColor: '#a8a8a8',
        text: 'item1'
    }
}

function createFormDeadline(taskName, deadline, stars) { // создание ui дедлайна
    const  div = document.createElement('div')
    div.style.cssText = `
      width: 200px;
      height: 100px;
      background-color: #353535;
      color: #FFFFFF
    `;
    div.innerHTML = `
        <h2>${taskName}</h2>
        <h3>${deadline}</h3>
    `;
    body.appendChild(div);
}

const buttonAddDeadline = document.getElementById('addDeadline');
buttonAddDeadline.addEventListener('click', function () // обработчик нажатия кнопки
{
    const formDataElem = document.getElementById('datetime-form');
    const taskName = document.getElementById('task-name').value;
    const deadline = document.getElementById('datetime-input').value
    createFormDeadline(taskName, deadline);

    fetch('/process', {method: 'POST', body: new FormData(formDataElem)}) // асинхронная отправка данных на go-сервер
        .then(response => response.text()) // вывод ответа от сервера
        .catch(error => console.error('Ошибка:', error)); // вывод ошибки, в случае ошибки
});