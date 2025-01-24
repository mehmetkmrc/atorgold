document.addEventListener("DOMContentLoaded", function () {
    const loginForm = document.querySelector('form[name="login-form"]');

    loginForm.addEventListener("submit", async function (event) {
        event.preventDefault(); // Formun otomatik olarak submit edilmesini engeller

        // Formdaki verileri almak
        const email = loginForm.querySelector('input[name="email"]').value;
        const password = loginForm.querySelector('input[name="password"]').value;

        // Veri objesi oluştur ve JSON formatına çevir
        const requestData = { email: email, password: password };

        try {
            const response = await fetch("/auth/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"  // JSON olarak gönderildiğini belirt
                },
                body: JSON.stringify(requestData)  // Objemizi JSON formatına çeviriyoruz
            });

            if (response.ok) {
                const data = await response.json();

                // Başarılı giriş
                console.log("Giriş başarılı:", data);
                window.location.href = "/add-product"; // Yönlendirme yap
            } else {
                // Hata mesajını göster
                const errorData = await response.json();
                console.error("Giriş hatası:", errorData.message);
                alert("Giriş başarısız. Lütfen bilgilerinizi kontrol edin.");
            }
        } catch (error) {
            console.error("Bir hata oluştu:", error);
            alert("Sunucuya bağlanırken bir sorun oluştu.");
        }
    });
});
