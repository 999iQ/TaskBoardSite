// кнопка создать аккаунт
const buttonCreateFormAccount = document.querySelector(".create-account-form");
buttonCreateFormAccount.addEventListener('click', async function (event) // обработчик нажатия кнопки
{
    event.preventDefault(); // для ссылки
    document.getElementById('login-form').hidden = true; //скрываем окно логина
    document.getElementById('reg-form').hidden = false; //показываем окно регистрации
});
// кнопка уже есть аккаунт
const buttonLoginFormAccount = document.querySelector(".i-have-account");
buttonLoginFormAccount.addEventListener('click', async function (event) // обработчик нажатия кнопки
{
    event.preventDefault(); // для ссылки
    document.getElementById('reg-form').hidden = true; //показываем окно регистрации
    document.getElementById('login-form').hidden = false; //скрываем окно логина
});
// кнопка войти
const buttonLoginAccount = document.getElementById("login-button");
buttonLoginAccount.addEventListener('click', async function () // обработчик нажатия кнопки
{
    const email = document.getElementById('email-log').value;
    const password = document.getElementById('password-log').value;
    // проверки на пустые поля
    if(email === "" || password === ""){
        alert("Поля для входа должны быть заполнены.");
        return;
    }
    // проверки на сходимость пароля и почты
    const formData = new FormData();
    formData.append("email", email);
    formData.append("password", password);

    const response = await fetch('/authorization',
        {method: 'POST', body: formData});

    console.log("login status POST /authorization: ", response.status)
    if (!response.ok) {
        if (response.status >= 300 && response.status < 400) { // редирект после входа
            // window.location.href = response.headers.get('Location'); // при обычной отправке статуса в Go, локация не передается
            window.location.href = "/"; // редирект
        }
        else {
            alert("Не верная почта или пароль.");
            throw new Error('login response was not ok');
        }
    }
    if (response.ok) {
        const responseData = await response.json();
        localStorage.setItem('token', responseData.token); // сохранение JWT токена в локальном хранилище у клиента
        window.location.href = "/"; // редирект
    }

});

// кнопка создать в форме регистрации аккаунта
const buttonRegAccount = document.getElementById("reg-button");
buttonRegAccount.addEventListener('click', async function () // обработчик нажатия кнопки
{
    const email = document.getElementById('email-reg').value;
    const nickname = document.getElementById('nickname').value;
    const password1 = document.getElementById('password-reg1').value;
    const password2 = document.getElementById('password-reg2').value;
    // проверки на пустые поля
    if(email === "" || password1 === "" || password2 === "" || nickname.length < 3){
        alert("Поля для входа должны быть заполнены. Длина ника больше 3 символов.");
        return;
    }
    if (password1 !== password2) {
        alert("Пароли не совпадают.");
        return;
    }
    // проверки на сходимость пароля и почты
    const formData = new FormData();
    formData.append("email", email);
    formData.append("nickname", nickname);
    formData.append("password", password1);

    const response = await fetch('/register',
    {method: 'POST', body: formData});

    if (!response.ok) {
        if (response.status >= 300 && response.status < 400) { // редирект после входа
            console.log("status", response.status)
            // window.location.href = response.headers.get('Location'); // при обычной отправке статуса в Go, локация не передается
            window.location.href = "/"; // редирект
        }
        else {
            alert("Аккаунт с такой почтой уже существует.");
            throw new Error('login response was not ok');
        }
    }
    if (response.ok) {
        const responseData = await response.json();
        localStorage.setItem('token', responseData.token); // сохранение JWT токена в локальном хранилище у клиента
        window.location.href = "/"; // редирект
    }
});