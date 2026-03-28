import { isUserLoggedIn, getPasswordPolicyError } from "./utils.js";
const $form = document.querySelector(".login-form");
const $registerBtn = document.querySelector(".register-btn");
const $emailInput = document.getElementById("email");
const $passwordInput = document.getElementById("password");
const $submitBtn = $form.querySelector("button[type='submit']");
const $formFeedback = document.getElementById("form-feedback");
const $emailError = document.getElementById("email-error");
const $passwordError = document.getElementById("password-error");

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

    if ($input.validity.typeMismatch) {
        setFieldError($input, $errorContainer, "Please enter a valid email address.");
        return false;
    }

    if ($input === $passwordInput) {
        const passwordPolicyError = getPasswordPolicyError($input.value);
        if (passwordPolicyError) {
            setFieldError($input, $errorContainer, passwordPolicyError);
            return false;
        }
    }

    return true;
}

$emailInput.addEventListener("input", () => {
    if ($emailInput.validity.valid) {
        clearFieldError($emailInput, $emailError);
    }
});

$passwordInput.addEventListener("input", () => {
    const passwordPolicyError = getPasswordPolicyError($passwordInput.value);
    if ($passwordInput.validity.valid && !passwordPolicyError) {
        clearFieldError($passwordInput, $passwordError);
    }
});

$form.addEventListener("submit", async function (event) {
    event.preventDefault();
    clearFormFeedback();

    const isEmailValid = validateField($emailInput, $emailError);
    const isPasswordValid = validateField($passwordInput, $passwordError);
    if (!isEmailValid || !isPasswordValid) {
        setFormFeedback("Please fix the highlighted fields and try again.");
        return;
    }

    const formData = new FormData($form);

    const email = formData.get("email");
    const password = formData.get("password");

    const loginReq = {
        "email": email,
        "password": password
    };

    try {
        $submitBtn.disabled = true;
        $submitBtn.textContent = "Signing in...";
        $form.setAttribute("aria-busy", "true");

        const url = "/api/v1/auth/login";
        const response = await fetch(url, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(loginReq)
        });

        if (!response.ok) {
            let message = "Email or password is incorrect.";
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

        window.location.href = "/";
    } catch (error) {
        setFormFeedback(error.message || "Unable to sign in right now. Please try again.");
    } finally {
        $submitBtn.disabled = false;
        $submitBtn.textContent = "Accept";
        $form.removeAttribute("aria-busy");
    }
});

$registerBtn.addEventListener("click", () => {
    window.location.href = "/register";
});

window.addEventListener("load", async () => {
    if (await isUserLoggedIn()) {
        window.location.href = "/";
    }
});