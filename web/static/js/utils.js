export async function isUserLoggedIn() {
    const url = "/api/v1/users/me"
    try {
        const response = await fetch(url, {method: "GET"})
        return response.ok
    } catch (error) {
        console.error(error.message)
        return false
    }
}

export async function logout() {
    if (!await isUserLoggedIn()) {
        console.error("invalid logout, log in")
        return
    }
    try {
        await fetch("/api/v1/auth/logout", {method: "DELETE"})
        window.location.reload();
    } catch (error) {
        console.error(error)
    }
}

export function renderComingSoonMovies($futureMovies, movies) {
    $futureMovies.innerHTML = '<h2>Coming Soon</h2>';

    if (!Array.isArray(movies) || movies.length === 0) {
        const emptyState = document.createElement('p');
        emptyState.textContent = 'No upcoming movies available right now.';
        $futureMovies.appendChild(emptyState);
        return;
    }

    const grid = document.createElement('div');
    grid.className = 'soon-movies-grid';

    for (const movie of movies) {
        const card = document.createElement('article');
        card.className = 'soon-movie-card';

        const poster = document.createElement('img');
        poster.src = movie.posterImageUrl;
        poster.alt = `${movie.title} poster`;
        poster.loading = 'lazy';

        const title = document.createElement('h3');
        title.className = 'soon-movie-title';
        title.textContent = movie.title;

        card.append(poster, title);
        card.style.cursor = "pointer";

        card.addEventListener("click", () => {
            window.location.href = `/movies/${movie.movieId}`
        })
        grid.appendChild(card);
    }

    $futureMovies.appendChild(grid);
}

export function getPasswordPolicyError(password) {
    if (password.length < 8 || password.length >= 32) {
        return "Password must be between 8 and 31 characters.";
    }

    const hasUpper = /\p{Lu}/u.test(password);
    const hasLower = /\p{Ll}/u.test(password);
    const hasDigit = /\p{Nd}/u.test(password);
    const hasSpecial = /[\p{P}\p{S}]/u.test(password);

    if (!hasUpper || !hasLower || !hasDigit || !hasSpecial) {
        return "Password must include uppercase, lowercase, number, and special character.";
    }

    return "";
}