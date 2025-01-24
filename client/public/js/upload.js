let imageBase64Array = [];

document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("documentLoader");
    const uploadInput = document.getElementById("upload");

    // Image yükleme ve base64 dönüşümü
    uploadInput.addEventListener('change', function() {
        const files = this.files;  // Birden fazla dosya alıyoruz
        const reader = new FileReader();

        // Her bir dosya için base64 verisini alalım
        Array.from(files).forEach(file => {
            reader.onload = function(event) {
                // Her resmin Base64 verisini diziye ekliyoruz
                imageBase64Array.push(event.target.result);
                displayImages(); // Resimleri görüntüle
            };

            if (file) {
                reader.readAsDataURL(file);
            } 
        });

        if (files.length > 0) {
            document.getElementById('upload-label').textContent = `${files.length} file(s) selected`; // Dosya seçildiğinde label'ı değiştir
        }
    });

      // Resimleri dinamik olarak görüntüleme
    function displayImages() {
        const imageArea = document.getElementById("imageArea");
        imageArea.innerHTML = ''; // Önceki resimleri temizle

        imageBase64Array.forEach((base64, index) => {
            // Resim div'i
            const imageWrapper = document.createElement("div");
            imageWrapper.classList.add("position-relative", "mx-2", "mb-3");

            // Resim etiketi
            const imgElement = document.createElement("img");
            imgElement.src = base64;
            imgElement.alt = "Uploaded Image";
            imgElement.classList.add("img-fluid", "rounded", "shadow-sm");

            // Silme butonu
            const removeButton = document.createElement("button");
            removeButton.textContent = "X";
            removeButton.classList.add("btn", "btn-danger", "btn-sm", "position-absolute", "top-0", "end-0", "translate-middle");
            removeButton.style.zIndex = "10";

            // Silme butonuna tıklama işlevi
            removeButton.addEventListener("click", function () {
                removeImage(index); // Resmi kaldır
            });

            // Resim ve butonu wrapper'a ekle
            imageWrapper.appendChild(imgElement);
            imageWrapper.appendChild(removeButton);

            // Wrapper'ı imageArea'ya ekle
            imageArea.appendChild(imageWrapper);
        });
    }

    // Belirli bir resmi kaldır
    function removeImage(index) {
        imageBase64Array.splice(index, 1); // Diziden resmi kaldır
        displayImages(); // Görüntülenen resimleri güncelle
    }

    form.addEventListener("submit", async function (event) {
        event.preventDefault();

        const loader = document.getElementById("loader");
        loader.classList.remove("d-none");

        try {
            // Step 1: Create Main Document
        const mainTitle = document.querySelector('input[name="main_title"]').value;
        const mainResponse = await fetch("http://127.0.0.1:3000/documenter/main", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ main_title: mainTitle }),
        });

        if (!mainResponse.ok) {
            showModal("error", "Hata!", "Ana belge oluşturulurken bir hata oluştu!");
            return;
        }

        const mainData = await mainResponse.json();

        // Step 2: Create Sub Document
        const subTitle = document.querySelector('input[name="sub_title"]').value;
        const productCode = document.querySelector('input[name="product_code"]').value;
        const subMessage = document.querySelector('textarea[name="sub_message"]').value;
        const colText = document.querySelector('textarea[name="about_collection"]').value;
        const jewCare = document.querySelector('textarea[name="jewellery_care"]').value;

        const subDocumentData = {
            main_id: mainData.data,
            sub_title: subTitle,
            product_code: productCode,
            sub_message: subMessage,
        };

        if (imageBase64Array.length > 0) {
            // Resimleri Base64 formatında array'e ekliyoruz
            subDocumentData.asset = imageBase64Array.map(base64 => base64.split(",")[1]);  // Base64 verilerini temizliyoruz
        }

        const subResponse = await fetch("http://127.0.0.1:3000/documenter/sub", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(subDocumentData),
        });

        if (!subResponse.ok) {
            const errorText = await subResponse.text();
            showModal("error", "Hata!", `Alt belge oluşturulurken bir hata oluştu! Hata: ${errorText}`);
            console.log("Sub Document Error: ", errorText);
            return;
        }

        const subData = await subResponse.json();
        const contentResponse = await fetch("http://127.0.0.1:3000/documenter/content", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                sub_id: subData.data,
                about_collection: colText,
                jewellery_care: jewCare,
            }),
        });

        if(!contentResponse.ok){
          showModal("error", "Hata!", `İçerik oluşturulurken bir hata oluştu!`);
          return;
        }

        showModal("success", "Başarılı!", "Ürün başarıyla eklendi!");
        } catch (error) {
            console.error("Hata oluştu: ", error);
            showModal("error", "Hata!", "Bir hata oluştu!");
        } finally {
            // Spinner'ı gizle
            loader.classList.add("d-none");
        }

        // Modal kapandıktan sonra sayfayı yenile
        const modalElement = document.getElementById("kt_modal_1");
        const modal = bootstrap.Modal.getInstance(modalElement);

        modalElement.addEventListener("hidden.bs.modal", function () {
            window.location.reload(); // Sayfayı yenile
        });

        modal.hide();
    });
});


function showModal(type, title, message) {
    const modalTitle = document.getElementById("kt_modal_1").querySelector(".modal-title");
    const modalBody = document.getElementById("kt_modal_1").querySelector(".modal-body");
    const modalFooter = document.getElementById("kt_modal_1").querySelector(".modal-footer");

    if (type === 'success') {
        modalTitle.textContent = title;
        modalBody.innerHTML = `<p class="text-success">${message}</p>`;
        modalFooter.innerHTML = `<button type="button" class="btn btn-light" data-bs-dismiss="modal">Kapat</button>`;
    } else if (type === 'error') {
        modalTitle.textContent = title;
        modalBody.innerHTML = `<p class="text-danger">${message}</p>`;
        modalFooter.innerHTML = `<button type="button" class="btn btn-light" data-bs-dismiss="modal">Kapat</button>`;
    }

    const modal = new bootstrap.Modal(document.getElementById("kt_modal_1"));
    modal.show();
}