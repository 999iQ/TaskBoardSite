alert("Это JavaScript!");

const button = document.getElementById('addDeadline');
const formElem = document.getElementById('datetime-form')
button.addEventListener('click', function () // обработчик нажатия кнопки
{
    fetch('/process', {method: 'POST', body: new FormData(formElem)}) // отправка данных пост запросом
        .then(response => response.text())
        .then(data => alert(data))
        .catch(error => console.error('Ошибка:', error));
});