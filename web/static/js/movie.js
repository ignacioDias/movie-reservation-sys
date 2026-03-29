import { isUserLoggedIn } from "./utils.js";

const $movieDetails = document.querySelector('.movie-details');
const $movieError = document.querySelector('.movie-error');
const $movieProjections = document.querySelector('.movie-projections');
const $currentMoviesBtn = document.querySelector('.current-movies-btn');
const $soonMoviesBtn = document.querySelector('.soon-movies-btn');
const $signUpBtn = document.querySelector('.signup-btn');
const $loginBtn = document.querySelector('.login-btn');
const $profileBtn = document.querySelector('.profile-btn');
const $logoutBtn = document.querySelector('.logout-btn');
const $loggedDivBtns = document.querySelector('.logged')
const $notLoggedDivBtns = document.querySelector('.not-logged')
const $goBackBtn = document.querySelector('.go-back-btn')

function getYoutubeEmbedUrl(trailerUrl) {
    if (!trailerUrl) return null;

    try {
        const parsedUrl = new URL(trailerUrl);
        let videoId = '';

        if (parsedUrl.hostname.includes('youtube.com')) {
            videoId = parsedUrl.searchParams.get('v') || '';
        } else if (parsedUrl.hostname.includes('youtu.be')) {
            videoId = parsedUrl.pathname.replace('/', '');
        }

        return videoId ? `https://www.youtube.com/embed/${videoId}` : null;
    } catch {
        return null;
    }
}

function formatReleaseDate(releaseDate) {
    if (!releaseDate) return 'TBA';

    const parsedDate = new Date(releaseDate);
    if (Number.isNaN(parsedDate.getTime())) return 'TBA';

    return parsedDate.toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

// Groups a flat projections array into { "YYYY-MM-DD": [projection, ...] }
function groupProjectionsByDate(projections) {
    const groups = {};
    for (const projection of projections) {
        const raw = projection.startsAt ?? projection.datetime ?? '';
        const date = raw.slice(0, 10);
        if (!date) continue;
        if (!groups[date]) groups[date] = [];
        groups[date].push(projection);
    }
    return groups;
}

function formatDateHeading(dateStr) {
    const date = new Date(dateStr + 'T00:00:00'); // force local timezone
    return date.toLocaleDateString(undefined, {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

function formatTime(projection) {
    const raw = projection.startsAt ?? projection.datetime ?? '';
    const date = new Date(raw);
    return date.toLocaleTimeString(undefined, {
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
    });
}

function renderProjections(projections) {
    $movieProjections.innerHTML = '<h2 class="projections-title">Showtimes</h2>';
    if (!Array.isArray(projections) || projections.length === 0) {
        const empty = document.createElement('p');
        empty.className = 'projections-empty';
        empty.textContent = 'No showtimes available at the moment.';
        $movieProjections.appendChild(empty);
        return;
    }

    const groups = groupProjectionsByDate(projections);

    for (const [date, dateProjections] of Object.entries(groups)) {
        const dateBlock = document.createElement('div');
        dateBlock.className = 'projection-date-block';

        const dateHeading = document.createElement('h3');
        dateHeading.className = 'projection-date';
        dateHeading.textContent = formatDateHeading(date);

        const timesRow = document.createElement('div');
        timesRow.className = 'projection-times';

        for (const projection of dateProjections) {
            const timeBtn = document.createElement('button');
            timeBtn.className = 'projection-time-btn';
            timeBtn.dataset.projectionId = projection.projectionId;

            const timeSpan = document.createElement('span');
            timeSpan.className = 'projection-time';
            timeSpan.textContent = formatTime(projection);

            const metaSpan = document.createElement('span');
            metaSpan.className = 'projection-meta';
            const parts = [projection.screeningFormat, projection.language].filter(Boolean);
            metaSpan.textContent = parts.join(' · ');

            timeBtn.append(timeSpan, metaSpan);

            timeBtn.addEventListener('click', () => {
                window.location.href = `/projections/${projection.projectionId}`;
            });

            timesRow.appendChild(timeBtn);
        }

        dateBlock.append(dateHeading, timesRow);
        $movieProjections.appendChild(dateBlock);
    }
}

async function displayProjections(movie) {
    try {
        const response = await fetch(`/api/v1/movies/id/${movie.movieId}/projections`);
        if (!response.ok) {
            throw new Error(`Error getting projections: ${response.status}`);
        }
        const result = await response.json();
        renderProjections(result);
    } catch (error) {
        console.error(error);
        const errMsg = document.createElement('p');
        errMsg.className = 'projections-empty';
        errMsg.textContent = 'Could not load showtimes.';
        $movieProjections.appendChild(errMsg);
    }
}

function renderMovie(movie) {
    const genres = Array.isArray(movie.genres) ? movie.genres : [];
    const embedUrl = getYoutubeEmbedUrl(movie.trailerUrl);
    const releaseDate = formatReleaseDate(movie.releaseDate);

    $movieDetails.innerHTML = '';

    const article = document.createElement('article');
    article.className = 'movie-card';

    const mediaSection = document.createElement('section');
    mediaSection.className = 'movie-media';

    const poster = document.createElement('img');
    poster.className = 'movie-poster';
    poster.src = movie.posterImageUrl;
    poster.alt = `${movie.title} poster`;
    poster.loading = 'lazy';

    mediaSection.appendChild(poster);

    const contentSection = document.createElement('section');
    contentSection.className = 'movie-content';

    const title = document.createElement('h1');
    title.className = 'movie-title';
    title.textContent = movie.title;
    document.title = movie.title;
    
    const metadata = document.createElement('div');
    metadata.className = 'movie-metadata';
    metadata.textContent = `Release date: ${releaseDate}`;

    const genresList = document.createElement('ul');
    genresList.className = 'movie-genres';

    for (const genre of genres) {
        const genreItem = document.createElement('li');
        genreItem.className = 'movie-genre';
        genreItem.textContent = genre;
        genresList.appendChild(genreItem);
    }

    const description = document.createElement('p');
    description.className = 'movie-description';
    description.textContent = movie.description || 'No description available.';

    contentSection.append(title, metadata, genresList, description);

    if (embedUrl) {
        const trailerSection = document.createElement('section');
        trailerSection.className = 'movie-trailer';

        const trailerTitle = document.createElement('h2');
        trailerTitle.textContent = 'Trailer';

        const iframe = document.createElement('iframe');
        iframe.src = embedUrl;
        iframe.title = `${movie.title} trailer`;
        iframe.allow = 'accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture';
        iframe.referrerPolicy = 'strict-origin-when-cross-origin';
        iframe.allowFullscreen = true;
        iframe.loading = 'lazy';

        trailerSection.append(trailerTitle, iframe);
        contentSection.appendChild(trailerSection);
    }

    article.append(mediaSection, contentSection);
    $movieDetails.appendChild(article);

    $movieProjections.style.display = 'block';
    displayProjections(movie);
}

window.addEventListener('load', async () => {
    if (await isUserLoggedIn()) {
        $loggedDivBtns.style.display = "block";
        $notLoggedDivBtns.style.display = "none";
    } 
    const movieId = window.location.pathname.split('/').pop();
    try {
        const response = await fetch(`/api/v1/movies/id/${movieId}`);
        if (!response.ok) {
            throw new Error(`Error getting movie: ${response.status}`);
        }
        const result = await response.json();
        renderMovie(result);
    } catch (error) {
        console.error(error.message);
        $movieError.textContent = 'Could not load movie details. Please try again later.';
        $movieError.style.display = 'block';
    }
});

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

$goBackBtn.addEventListener("click", () => {
    window.location.href = "/";
})