
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

    const response = await fetch('/api/authorization',
        {method: 'POST', body: formData});

    if (!response.ok) {
        if (response.status >= 300 && response.status < 400) { // редирект после входа
            console.log("status", response.status)
            // console.log("location", response.headers.get('Location'))
            // window.location.href = response.headers.get('Location'); // при обычной отправке статуса в Go, локация не передается
            window.location.href = "/"; // редирект
        }
        else {
            throw new Error('login response was not ok');
        }
    }

});

// кнопка создать в форме регистрацииаккаунт
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

    const response = await fetch('/api/register',
{method: 'POST', body: formData});
    if (!response.ok) {
        throw new Error('getTasks response was not ok');
    }

    // response.text().then(result => {
    //     console.log("result: ", result);
    //     if(result === "failed") { // успешный логин
    //         console.log("Вход не выполнен");
    //     }
        // else { // logged in success and render page
        //     document.body.innerHTML = html;
        // }
    // });
});
