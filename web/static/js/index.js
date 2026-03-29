import { isUserLoggedIn, logout, renderComingSoonMovies } from "./utils.js";


const $currentMoviesBtn = document.querySelector('.current-movies-btn');
const $soonMoviesBtn = document.querySelector('.soon-movies-btn');
const $signUpBtn = document.querySelector('.signup-btn');
const $loginBtn = document.querySelector('.login-btn');
const $profileBtn = document.querySelector('.profile-btn');
const $logoutBtn = document.querySelector('.logout-btn');
const $loggedDivBtns = document.querySelector('.logged')
const $notLoggedDivBtns = document.querySelector('.not-logged')
const $currentMovies = document.querySelector('.current-projections')
const $futureMovies = document.querySelector('.future-movies')

function createMovieCard(movie) {
    const card = document.createElement('article');
    card.className = 'movie-card';

    const poster = document.createElement('img');
    poster.src = movie.posterImageUrl;
    poster.alt = `${movie.title} poster`;
    poster.loading = 'lazy';

    const title = document.createElement('h3');
    title.className = 'movie-title';
    title.textContent = movie.title;
    card.style.cursor = "pointer";
    card.addEventListener("click", () => {
        window.location.href = `/movies/id/${movie.movieId}`
    })
    card.append(poster, title);
    return card;
}

function renderCurrentMovies(movies) {
    $currentMovies.innerHTML = '<h2>Now Showing</h2>';

    if (!Array.isArray(movies) || movies.length === 0) {
        const emptyState = document.createElement('p');
        emptyState.textContent = 'No current movies available right now.';
        $currentMovies.appendChild(emptyState);
        return;
    }

    const carousel = document.createElement('div');
    carousel.className = 'movie-carousel';

    const track = document.createElement('div');
    track.className = 'movie-track';

    const doubledMovies = [...movies, ...movies];
    for (const movie of doubledMovies) {
        track.appendChild(createMovieCard(movie));
    }

    carousel.appendChild(track);
    $currentMovies.appendChild(carousel);
}




window.addEventListener('load', async (event) => {
    if (await isUserLoggedIn()) {
        $loggedDivBtns.style.display = "block";
        $notLoggedDivBtns.style.display = "none";
    } 
    try {
        const response = await fetch("/api/v1/movies/available_now");
        if(!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
        const result = await response.json();
        renderCurrentMovies(result);
    } catch (error) {
        console.error(error)
        $currentMovies.innerHTML = '<h2>Now Showing</h2><p>Could not load current movies.</p>';
    }
    try {
        const response = await fetch("/api/v1/movies/soon")
        if(!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
        const result = await response.json();
        renderComingSoonMovies($futureMovies, result)
    } catch (error) {
        console.error(error)
        $futureMovies.innerHTML = '<h2>Coming Soon</h2><p>Could not load upcoming movies.</p>';
    }
})

$currentMoviesBtn.addEventListener("click", () => {
    window.location.href = "/movies/projecting";
})
$soonMoviesBtn.addEventListener("click", () => {
    window.location.href = "/movies/soon";
})
$signUpBtn.addEventListener("click", () => {
    window.location.href = "/register";
})
$loginBtn.addEventListener("click", () => {
    window.location.href = "/login";
})
$logoutBtn.addEventListener("click", async () => {
    await logout();
})
$profileBtn.addEventListener("click", () => {
    window.location.href = "/me";
})