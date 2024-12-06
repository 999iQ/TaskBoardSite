const body = document.querySelector('body') // тело страницы для добавления туда новых форм дедлайнов

window.onload = async function() { //срабатывает при загрузке страницы
    await getTasks(); // заполнение списка задач
    console.log(new Date())
};

async function startTimer(duration, elementTimer) {
    let startTime = performance.now();
    function animate(currentTime) {
        const elapsedTime = (currentTime - startTime) / 1000 + 3*60*60; // вычитаем время загрузки страницы, а также корректируем часовой пояс (+3 часа)
        const remainingTime = duration - elapsedTime;

        // Math.abs - для отрицательного просроченного времени дедлайна
        const days = Math.floor(Math.abs(remainingTime / (60 * 60 * 24)));
        const hours = Math.floor(Math.abs((remainingTime % (60 * 60 * 24)) / (60 * 60)));
        const minutes = Math.floor(Math.abs((remainingTime % (60 * 60)) / 60));
        const seconds = Math.floor(Math.abs(remainingTime % 60));
        // Миллисекунды не добавляем для упрощения, но можно добавить аналогично предыдущему примеру

        const formattedDays = String(Math.abs(days)).padStart(2, '0');
        const formattedHours = String(Math.abs(hours)).padStart(2, '0');
        const formattedMinutes = String(Math.abs(minutes)).padStart(2, '0');
        const formattedSeconds = String(Math.abs(seconds)).padStart(2, '0');

        if (remainingTime < 0) { // добавляем текст "просрочено" для просроченных дедлайнов
            elementTimer.children[0].hidden = false;
        }
        elementTimer.children[2].textContent = `${formattedDays}`;
        elementTimer.children[3].textContent = `${formattedHours}`;
        elementTimer.children[4].textContent = `${formattedMinutes}`;
        elementTimer.children[5].textContent = `${formattedSeconds}`;
        requestAnimationFrame(animate);
    }
    requestAnimationFrame(animate);
}

async function getRemainingTimePercentage(startTimestamp, endTimestamp) {
    // функция считает оставшееся время до дедлайна в ПРОЦЕНТАХ для окраски таймера цветом
    // Приобразуем даты в милисекунды
    startTimestamp = new Date(startTimestamp).getTime() - 3 * 60 * 60 * 1000; // время в мили секах с поправкой на часовой пояс
    endTimestamp = new Date(endTimestamp).getTime() - 3 * 60 * 60 * 1000;

    // Вычисляем общую разницу во времени
    const totalTimeDifference = endTimestamp - startTimestamp;

    // Вычисляем прошедшее время
    const elapsedTime = Date.now() - startTimestamp;

    // Обработка случаев, когда endDate < startDate или endDate < now
    if (totalTimeDifference <= 0 || elapsedTime > totalTimeDifference) {
        return 0; // 100% если дата окончания прошла, или дата окончания раньше чем начала
    }
    // Вычисляем процент оставшегося времени
    return ((totalTimeDifference - elapsedTime) / totalTimeDifference) * 100; // type number
}

async function clearTasks() { // функция отчистки вводимых полей для кнопки отмена
    const tasks = document.querySelectorAll('.task');
    tasks.forEach(task => {
        task.remove();
    });
}

async function getTasks() { // заполнение списка задач
    try {
        const response = await fetch('/api/getTasks/0');
        if (!response.ok) { throw new Error('getTasks response was not ok'); }

        const tasks = await response.json(); // парс json'а
        await clearTasks() // отчистка списка задач перед заполнением
        // обработка данных => создание форм задач
        tasks.forEach(task => {
            createFormDeadline(task.Title, task.Deadlinedate, task.Priority, task.Status, task.Task_id) // создание форм списка задач

            const elementTimer = document.getElementById(`timer-${task.Task_id}`); // поиск таймера с таким же индексом
            if (task.Status === true){ // убираем таймер и подпись "просрочено" для выполненных задач
                elementTimer.children[0].hidden = true;
                elementTimer.children[2].hidden = true;
                elementTimer.children[3].hidden = true;
                elementTimer.children[4].hidden = true;
                elementTimer.children[5].hidden = true;
            } else {
                const duration = getSecondsUntilDate(task.Deadlinedate); // получаем разницу дат в секундах
                startTimer(duration, elementTimer); // запуск таймера дедлайна

                // работа над вычеслением процента оставщегося времени до дедлайна, с целью окрашивания таймера в нужный цвет
                getRemainingTimePercentage(task.Createdat, task.Deadlinedate).then(procent => {
                    if (procent <= 10) {
                        elementTimer.style.color = "#fa0005";
                    } else if (procent <= 30) {
                        elementTimer.style.color = "#f75e01";
                    } else if (procent <= 50) {
                        elementTimer.style.color = "#fbd101";
                    } else {
                        elementTimer.style.color = "#01a261";
                    }
                });
            }
        });
    } catch (error) {
        console.error('Error:', error);
    }
}

function getSecondsUntilDate(dateString) {
    const targetDate = new Date(dateString); // Парсинг строки в объект Date
    const now = new Date(); // текущее время с учетом часовых поясов
    const differenceInMilliseconds = targetDate - now; // Разница в миллисекундах
    return Math.floor(differenceInMilliseconds / 1000); // Разница в секундах
}

async function createFormDeadline(taskName, deadline, priority, status, taskId) { // создание ui дедлайна
    const  div = document.createElement('div')
    div.className = 'task'; // стили
    let stars = ""
    if (priority > 0) { // проставляем звездочки количеством в приоритет
        stars = "★".repeat(priority);
    }
    let statusCheckbox = "";
    if (status === true) {
        statusCheckbox = "checked";
    }

    div.innerHTML = `
        <input class="task-check-box" type="checkbox" id="checkbox-${taskId}" ${statusCheckbox}>
        <a href="" class="clickable-text" id="clickable-text-${taskId}">
            <span class="task-text">${taskName}</span>
        </a>
        <div class="task-timer" id="timer-${taskId}">
            <span class="task-text-delay" hidden>просрочено</span> 
            <span class="task-stars">${stars}</span> 
            <span class="task-time-box" id="days"></span>
            <span class="task-time-box" id="hours"></span>
            <span class="task-time-box" id="minutes"></span>
            <span class="task-time-box" id="seconds"></span>
        </div>
    `;
    const content_body = document.querySelector('.main-content');
    await content_body.appendChild(div);

    // ставим обработчики нажатий на название таски и на чекбокс
    const linkTaskTitle = div.querySelector(`#clickable-text-${taskId}`);
    await setClickableListener(linkTaskTitle, linkTaskTitle.id.split("-")[2]);

    const linkTaskCheckbox = div.querySelector(`#checkbox-${taskId}`);
    await setClickableListener(linkTaskCheckbox, linkTaskCheckbox.id.split("-")[1]);
}

async function validityInputForms(){
    // проверяем введённые значения в формах на корректность
    const taskName = document.getElementById('task-name').value;
    const description = document.getElementById('task-description').value;
    const deadline = document.getElementById('datetime-input').value;
    let valid = true;
    if (!/^[A-Za-zА-Яа-я0-9\s]{3,255}$/.test(taskName)) { // проверка на валидность имени таски
        alert("Невалидное название задачи.\nТолько буквы и цифры, длина от 3 до 255 символов");
        valid = false;
    }
    else if (description.length > 2000) { // проверка на валидность имени таски
        alert("Невалидное описание задачи.\nДлина до 2000 символов");
        valid = false;
    }
    else if (deadline === '') { // проверка на актуальность даты дедлайна
        alert("Невалидная дата дедлайна.");
        valid = false;
    }
    return valid;
}

const buttonCancelDeadline = document.getElementById('cancelDeadline');
buttonCancelDeadline.addEventListener('click', function ()
{
    clearInputForms();
    document.getElementsByClassName('sidebar')[0].hidden = true; // скрывает меню
})

function clearInputForms() {
    // Отчистка (сброс) всех форм ввода в боковом меню
    const forms = document.querySelectorAll('form');
    // Очищаем каждую форму
    forms.forEach(form => {
        form.reset();
    });
}

async function setClickableListener(link, id_number) { // функция установки обработчика нажатий на таске
    // активация меню по нажатию на название таски
    link.addEventListener('click', async function (event) {
        if (this.type !== "checkbox") {
            event.preventDefault(); // Предотвращаем стандартное поведение ссылки (только для кликабельного текста, игнорируем чекбокс, а то он не будет нажиматься =)
        }
        try {
            document.querySelectorAll('a').forEach(otherLink => {
                otherLink.style.color = '#FFFFFF'; // неактивным таскам возвращаем белый цвет
            });
            this.style.color = '#01a361'; // выделяю цветом выбранную таску

            document.getElementsByClassName('sidebar')[0].hidden = false; // показывает меню

            const response = await fetch(`/api/getTasks/${id_number}`);
            if (!response.ok) {
                throw new Error('getTasks response was not ok');
            }
            const task = await response.json(); // парс json'а
            document.getElementById('datetime-form').dataset.id = task[0].Task_id;

            // заполняем формы, правого окна создания дедлайнов
            document.getElementById('task-name').value = task[0].Title;
            document.getElementById('task-description').value = task[0].Description;
            document.getElementById('datetime-input').value = new Date(task[0].Deadlinedate).toISOString().slice(0, -5);
            document.getElementById(task[0].Priority).checked = true;

            // скрываем кнопки для добавления дедлайна
            document.getElementById('cancelDeadline').hidden = true;
            document.getElementById('addDeadline').hidden = true;
            // показываем кнопки удаления и изменения
            document.getElementById('deleteDeadline').hidden = false;
            document.getElementById('editDeadline').hidden = false;

            if (this.type === "checkbox") {
                await addAndEditTask(); // т.е. если нажали на чекбокс, то помимо активации меню сразу обновляем статус задачи в бд
            }
        } catch (error) {
            console.error('Error:', error);
        }
    })
}

document.addEventListener('click', function (event) {
    // Проверяем, был ли клик за пределами бокового меню, чтобы его скрыть
    const targetElementMenu = document.getElementsByClassName('sidebar')[0];
    const targetElementTasks = document.getElementsByClassName('task')
    const clickInsideTasks = Array.from(targetElementTasks).some(taskElement => taskElement.contains(event.target));
    const targetElementCreateButton = document.getElementById('createDeadline');

    if (!targetElementMenu.contains(event.target) && !clickInsideTasks && !targetElementCreateButton.contains(event.target)) { // клик в пустом месте
        document.getElementsByClassName('sidebar')[0].hidden = true; // скрывает меню
    }
});

const buttonCreateDeadline = document.getElementById('createDeadline');
buttonCreateDeadline.addEventListener('click', function () // обработчик нажатия кнопки
{
    clearInputForms();
    document.getElementById('datetime-form').dataset.id = '';
    document.getElementsByClassName('sidebar')[0].hidden = false; // показывает меню
    // скрываем кнопки удаления и изменения
    document.getElementById('deleteDeadline').hidden = true;
    document.getElementById('editDeadline').hidden = true;
    // показываем кнопки для добавления дедлайна
    document.getElementById('cancelDeadline').hidden = false;
    document.getElementById('addDeadline').hidden = false;
});

const buttonDeleteDeadline = document.getElementById('deleteDeadline');
buttonDeleteDeadline.addEventListener('click', async function () // обработчик нажатия кнопки
{
    // удаление таски через api
    const task_id = document.getElementById('datetime-form').dataset.id;
    clearInputForms();
    document.getElementsByClassName('sidebar')[0].hidden = true; // показывает меню
    await fetch(`/api/deleteDeadline/${task_id}`, {method: 'POST'});
    await getTasks();
});


document.getElementById('addDeadline').addEventListener('click', addAndEditTask); // обработчики кнопок: добавить и изменить вызывают один метод
document.getElementById('editDeadline').addEventListener('click', addAndEditTask); // обработчики кнопок: добавить и изменить вызывают один метод

async function addAndEditTask() {
    // изменение и добавление таски
    if (!await validityInputForms()) return;

    const taskName = document.getElementById('task-name').value;
    const description = document.getElementById('task-description').value;
    const deadline = document.getElementById('datetime-input').value;
    const taskId = document.getElementById('datetime-form').dataset.id;
    const status = taskId ? document.getElementById(`checkbox-${taskId}`).checked : false; // если создаём таску, то поле id в форме пустое

    let priority = document.querySelector('input[name="slider"]:checked'); // поиск активного radio-элемента (текущий приоритет задачи)
    priority = priority ? parseInt(priority.id) : 0;

    const formData = new FormData();
    formData.append("task-name", taskName);
    formData.append("task-description", description);
    formData.append("datetime-input", deadline);
    formData.append("priority", priority);
    formData.append("task-id", taskId);
    formData.append("status", status);

    await fetch('/api/addAndEditDeadline', {method: 'POST', body: formData}) // асинхронная отправка данных на go-сервер
        .then(response => response.text()) // вывод ответа от сервера
        .catch(error => console.error('Ошибка:', error)); // вывод ошибки, в случае ошибки

    await getTasks(); // отчистка и отрисовка ui тасок
}