import { getPasswordPolicyError, isUserLoggedIn } from "./utils.js";

const $form = document.querySelector(".register-form");
const $loginBtn = document.querySelector(".login-btn");
const $emailInput = document.getElementById("email");
const $documentNumberInput = document.getElementById("documentNumber");
const $passwordInput = document.getElementById("password");
const $passwordConfirmInput = document.getElementById("passwordConfirm");
const $profilePictureInput = document.getElementById("profilePicture");
const $avatarImages = document.querySelectorAll(".profile-pictures img");
const $submitBtn = $form.querySelector("button[type='submit']");
const $formFeedback = document.getElementById("form-feedback");
const $emailError = document.getElementById("email-error");
const $documentNumberError = document.getElementById("documentNumber-error");
const $passwordError = document.getElementById("password-error");
const $passwordConfirmError = document.getElementById("passwordConfirm-error");
const $avatarError = document.getElementById("avatar-error");

let selectedAvatarElement = null;

function setFieldError($input, $errorContainer, message) {
    $input.setAttribute("aria-invalid", "true");
    $errorContainer.textContent = message;
    $errorContainer.hidden = false;
}

function clearFieldError($input, $errorContainer) {
    $input.setAttribute("aria-invalid", "false");
    $errorContainer.textContent = "";
    $errorContainer.hidden = true;
}

function setFormFeedback(message) {
    $formFeedback.textContent = message;
    $formFeedback.hidden = false;
}

function clearFormFeedback() {
    $formFeedback.textContent = "";
    $formFeedback.hidden = true;
}

function validateField($input, $errorContainer) {
    clearFieldError($input, $errorContainer);

    if ($input.validity.valueMissing) {
        setFieldError($input, $errorContainer, "This field is required.");
        return false;
    }

    if ($input === $emailInput && $input.validity.typeMismatch) {
        setFieldError($input, $errorContainer, "Please enter a valid email address.");
        return false;
    }

    if ($input === $documentNumberInput && $input.value.trim() === "") {
        setFieldError($input, $errorContainer, "Document number cannot be empty.");
        return false;
    }

    if ($input === $passwordInput) {
        const passwordPolicyError = getPasswordPolicyError($input.value);
        if (passwordPolicyError) {
            setFieldError($input, $errorContainer, passwordPolicyError);
            return false;
        }
    }

    if ($input === $passwordConfirmInput) {
        if ($input.value !== $passwordInput.value) {
            setFieldError($input, $errorContainer, "Passwords do not match.");
            return false;
        }
    }

    return true;
}

$avatarImages.forEach(img => {
    img.addEventListener("click", () => {
        if (selectedAvatarElement) {
            selectedAvatarElement.style.border = "none";
        }
        selectedAvatarElement = img;
        img.style.border = "2px solid #c8a96e";
        $profilePictureInput.value = img.dataset.avatar;
        clearFieldError($profilePictureInput, $avatarError);
    });

    img.addEventListener("keypress", (event) => {
        if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            img.click();
        }
    });
});

$emailInput.addEventListener("input", () => {
    if ($emailInput.validity.valid) {
        clearFieldError($emailInput, $emailError);
    }
});

$documentNumberInput.addEventListener("input", () => {
    if ($documentNumberInput.value.trim() !== "") {
        clearFieldError($documentNumberInput, $documentNumberError);
    }
});

$passwordInput.addEventListener("input", () => {
    const passwordPolicyError = getPasswordPolicyError($passwordInput.value);
    if ($passwordInput.validity.valid && !passwordPolicyError) {
        clearFieldError($passwordInput, $passwordError);
    }
    if ($passwordConfirmInput.value) {
        if ($passwordConfirmInput.value === $passwordInput.value) {
            clearFieldError($passwordConfirmInput, $passwordConfirmError);
        }
    }
});

$passwordConfirmInput.addEventListener("input", () => {
    if ($passwordConfirmInput.value === $passwordInput.value) {
        clearFieldError($passwordConfirmInput, $passwordConfirmError);
    }
});

$form.addEventListener("submit", async function (event) {
    event.preventDefault();
    clearFormFeedback();

    const isEmailValid = validateField($emailInput, $emailError);
    const isDocumentNumberValid = validateField($documentNumberInput, $documentNumberError);
    const isPasswordValid = validateField($passwordInput, $passwordError);
    const isPasswordConfirmValid = validateField($passwordConfirmInput, $passwordConfirmError);

    if (!$profilePictureInput.value) {
        setFieldError($profilePictureInput, $avatarError, "Please select an avatar.");
    }

    if (!isEmailValid || !isDocumentNumberValid || !isPasswordValid || !isPasswordConfirmValid || !$profilePictureInput.value) {
        setFormFeedback("Please fix the highlighted fields and try again.");
        return;
    }

    const registerReq = {
        "email": $emailInput.value,
        "password": $passwordInput.value,
        "documentNumber": $documentNumberInput.value,
        "profilePicture": $profilePictureInput.value
    };

    try {
        $submitBtn.disabled = true;
        $submitBtn.textContent = "Creating account...";
        $form.setAttribute("aria-busy", "true");

        const url = "/api/v1/auth/register";
        const response = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(registerReq)
        });
        console.log(JSON.stringify(registerReq))
        if (!response.ok) {
            let message = "Unable to register. Please check your information and try again.";
            try {
                const data = await response.json();
                if (typeof data?.message === "string" && data.message.trim()) {
                    message = data.message;
                }
            } catch {
                // Keep fallback message when response body is not JSON.
            }
            throw new Error(message);
        }

        window.location.href = "/login";
    } catch (error) {
        setFormFeedback(error.message || "Unable to register right now. Please try again.");
    } finally {
        $submitBtn.disabled = false;
        $submitBtn.textContent = "Register";
        $form.removeAttribute("aria-busy");
    }
});

$loginBtn.addEventListener("click", () => {
    window.location.href = "/login";
});

window.addEventListener("load", async () => {
    if(await isUserLoggedIn()) {
        window.location.href = "/"
    }
})
