import { isUserLoggedIn, logout } from "./utils.js";

// ── DOM refs ──────────────────────────────────────────────────────────────────
const $currentMoviesBtn  = document.querySelector('.current-movies-btn');
const $soonMoviesBtn     = document.querySelector('.soon-movies-btn');
const $signUpBtn         = document.querySelector('.signup-btn');
const $loginBtn          = document.querySelector('.login-btn');
const $profileBtn        = document.querySelector('.profile-btn');
const $logoutBtn         = document.querySelector('.logout-btn');
const $loggedDivBtns     = document.querySelector('.logged');
const $notLoggedDivBtns  = document.querySelector('.not-logged');
const $goBackBtn         = document.querySelector('.go-back-btn');
const $tickets           = document.querySelector('.tickets');
const $confirmBtn        = document.querySelector('.confirm');

// ── State ─────────────────────────────────────────────────────────────────────
let availableTickets = [];          // [{ticketId, name, price, cant_seats}, ...]
let ticketQuantities = {};          // { ticketId: quantity }
let selectedSeats    = new Set();   // Set of "row-col" strings
let auditoriumRows   = 8;
let auditoriumCols   = 10;

// ── Helpers ───────────────────────────────────────────────────────────────────
function getProjectionIdFromURL() {
    // URL pattern: /projections/id/{projection_id}/reservation
    const parts = window.location.pathname.split('/');
    const idx = parts.indexOf('id');
    return idx !== -1 ? parts[idx + 1] : null;
}

function totalSeatsSelected() {
    return Object.values(ticketQuantities).reduce((sum, q) => sum + q, 0);
}

function totalPrice() {
    return availableTickets.reduce((sum, t) => {
        return sum + (ticketQuantities[t.ticketId] || 0) * t.price;
    }, 0);
}

// ── Fetch tickets ─────────────────────────────────────────────────────────────
async function loadTickets() {
    try {
        const res = await fetch('/api/v1/tickets', { credentials: 'include' });
        if (!res.ok) throw new Error('Failed to fetch tickets');
        availableTickets = await res.json();
        availableTickets.forEach(t => { ticketQuantities[t.ticketId] = 0; });
        renderTickets();
    } catch (err) {
        showError('Could not load ticket types. Please try again.');
        console.error(err);
    }
}

// ── Fetch projection (to get auditorium size) ─────────────────────────────────
async function loadProjection(projectionId) {
    if (!projectionId) return;
    try {
        const res = await fetch(`/api/v1/projections/${projectionId}`, { credentials: 'include' });
        if (!res.ok) return;
        const projection = await res.json();
        if (projection.auditoriumId) {
            const audRes = await fetch(`/api/v1/auditoriums/${projection.auditoriumId}`, { credentials: 'include' });
            if (audRes.ok) {
                const aud = await audRes.json();
                auditoriumRows = aud.cantRows || 8;
                auditoriumCols = aud.cantCols || 10;
                renderSeatMap();
            }
        }
    } catch (err) {
        console.warn('Could not load projection details:', err);
    }
}

// ── Render ────────────────────────────────────────────────────────────────────
function renderAll() {
    renderTickets();
    renderSeatMap();
    renderSummary();
}

function renderTickets() {
    let html = `<div class="section-title">Select Tickets</div>
                <div class="ticket-list">`;

    if (availableTickets.length === 0) {
        html += `<p class="no-tickets">No ticket types available.</p>`;
    }

    availableTickets.forEach(t => {
        const qty = ticketQuantities[t.ticketId] || 0;
        html += `
        <div class="ticket-card">
            <div class="ticket-info">
                <span class="ticket-name">${t.name}</span>
                <span class="ticket-price">$${t.price.toFixed(2)}</span>
            </div>
            <div class="ticket-counter">
                <button class="counter-btn minus" data-id="${t.ticketId}" data-step="${t.cant_seats}" ${qty === 0 ? 'disabled' : ''}>−</button>
                <span class="counter-value">${qty}</span>
                <button class="counter-btn plus" data-id="${t.ticketId}" data-step="${t.cant_seats}">+</button>
            </div>
        </div>`;
    });

    html += `</div>`;

    // Summary bar
    const total = totalPrice();
    const seatsNeeded = totalSeatsSelected();
    html += `
    <div class="summary-bar">
        <div class="summary-seats">
            <span class="summary-label">Seats needed</span>
            <span class="summary-value">${seatsNeeded}</span>
        </div>
        <div class="summary-total">
            <span class="summary-label">Total</span>
            <span class="summary-price">$${total.toFixed(2)}</span>
        </div>
    </div>`;

    // Inject ticket section only
    let ticketSection = document.querySelector('.ticket-section');
    if (!ticketSection) {
        // First render — build full layout
        $tickets.innerHTML = `
            <div class="ticket-section">${html}</div>
            <div class="section-title seat-section-title">Choose Your Seats</div>
            <div class="seat-section"></div>
            <div class="error-msg"></div>
        `;
        renderSeatMap();
    } else {
        ticketSection.innerHTML = html;
    }

    // Re-query confirm button after innerHTML replacement and append it
    const existingConfirm = document.querySelector('.tickets > .confirm');
    if (!existingConfirm) {
        $tickets.appendChild($confirmBtn);
    }

    attachTicketListeners();
    updateConfirmState();
}

function attachTicketListeners() {
    document.querySelectorAll('.counter-btn.plus').forEach(btn => {
        btn.addEventListener('click', () => {
            const id   = parseInt(btn.dataset.id);
            const step = parseInt(btn.dataset.step) || 1;
            ticketQuantities[id] = (ticketQuantities[id] || 0) + step;
            syncSeatsToQuota();
            renderTickets();
            renderSeatMap();
        });
    });
    document.querySelectorAll('.counter-btn.minus').forEach(btn => {
        btn.addEventListener('click', () => {
            const id   = parseInt(btn.dataset.id);
            const step = parseInt(btn.dataset.step) || 1;
            if ((ticketQuantities[id] || 0) >= step) {
                ticketQuantities[id] -= step;
                syncSeatsToQuota();
                renderTickets();
                renderSeatMap();
            }
        });
    });
}

function syncSeatsToQuota() {
    const needed = totalSeatsSelected();
    if (selectedSeats.size > needed) {
        const arr = [...selectedSeats];
        selectedSeats = new Set(arr.slice(0, needed));
    }
}

function renderSeatMap() {
    const $seatSection = document.querySelector('.seat-section');
    if (!$seatSection) return;

    const needed = totalSeatsSelected();
    const chosen = selectedSeats.size;

    let html = `
    <div class="screen-label">SCREEN</div>
    <div class="screen-bar"></div>
    <div class="seat-grid" style="--cols:${auditoriumCols}">`;

    for (let row = 1; row <= auditoriumRows; row++) {
        html += `<div class="seat-row"><span class="row-label">${String.fromCharCode(64 + row)}</span>`;
        for (let col = 1; col <= auditoriumCols; col++) {
            const key = `${row}-${col}`;
            const isSelected = selectedSeats.has(key);
            const isDisabled = !isSelected && chosen >= needed && needed > 0;
            html += `<button 
                class="seat ${isSelected ? 'selected' : ''} ${isDisabled ? 'disabled' : ''}"
                data-row="${row}" 
                data-col="${col}"
                ${isDisabled ? 'disabled' : ''}
                title="Row ${String.fromCharCode(64+row)}, Seat ${col}"
            ></button>`;
        }
        html += `</div>`;
    }

    html += `</div>
    <div class="seat-legend">
        <span class="legend-item"><span class="legend-dot available"></span> Available</span>
        <span class="legend-item"><span class="legend-dot selected"></span> Selected (${chosen}/${needed})</span>
        <span class="legend-item"><span class="legend-dot taken"></span> Unavailable</span>
    </div>`;

    $seatSection.innerHTML = html;

    $seatSection.querySelectorAll('.seat:not([disabled])').forEach(btn => {
        btn.addEventListener('click', () => {
            const key = `${btn.dataset.row}-${btn.dataset.col}`;
            if (selectedSeats.has(key)) {
                selectedSeats.delete(key);
            } else if (selectedSeats.size < needed) {
                selectedSeats.add(key);
            }
            renderSeatMap();
            updateConfirmState();
        });
    });

    updateConfirmState();
}

function updateConfirmState() {
    const needed  = totalSeatsSelected();
    const chosen  = selectedSeats.size;
    const ready   = needed > 0 && chosen === needed;
    $confirmBtn.disabled = !ready;
    $confirmBtn.classList.toggle('ready', ready);
}

function showError(msg) {
    const $err = document.querySelector('.error-msg');
    if ($err) {
        $err.textContent = msg;
        $err.style.display = 'block';
        setTimeout(() => { $err.style.display = 'none'; }, 4000);
    }
}

// ── Submit reservation ────────────────────────────────────────────────────────
$confirmBtn.addEventListener('click', async () => {
    const projectionId = getProjectionIdFromURL();
    if (!projectionId) {
        showError('Missing projection ID in URL.');
        return;
    }

    const seats = [...selectedSeats].map(key => {
        const [row, col] = key.split('-').map(Number);
        return { row, col };
    });

    const tickets = availableTickets
        .filter(t => (ticketQuantities[t.ticketId] || 0) > 0)
        .map(t => ({
            ticketId:   t.ticketId,
            name:       t.name,
            price:      t.price,
            cant_seats: ticketQuantities[t.ticketId]
        }));

    const body = {
        projectionId: parseInt(projectionId),
        seats,
        tickets
    };

    $confirmBtn.disabled = true;
    $confirmBtn.textContent = 'Confirming…';

    try {
        const res = await fetch('/api/v1/reservations', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });

        if (!res.ok) {
            const errText = await res.text();
            throw new Error(errText || 'Reservation failed');
        }

        const reservation = await res.json();
        window.location.href = `/`;
    } catch (err) {
        showError(err.message || 'Could not complete reservation. Try again.');
        $confirmBtn.disabled = false;
        $confirmBtn.textContent = 'Accept';
        updateConfirmState();
    }
});

// ── Nav listeners ─────────────────────────────────────────────────────────────
$currentMoviesBtn.addEventListener('click', () => { window.location.href = '/movies/projecting'; });
$soonMoviesBtn.addEventListener('click',    () => { window.location.href = '/movies/soon'; });
$signUpBtn.addEventListener('click',        () => { window.location.href = '/register'; });
$loginBtn.addEventListener('click',         () => { window.location.href = '/login'; });
$logoutBtn.addEventListener('click', async () => { await logout(); });
$profileBtn.addEventListener('click',       () => { window.location.href = '/me'; });
$goBackBtn.addEventListener('click',        () => { window.history.back(); });

// ── Boot ──────────────────────────────────────────────────────────────────────
window.addEventListener('load', async () => {
    if (await isUserLoggedIn()) {
        $loggedDivBtns.style.display = 'block';
        $notLoggedDivBtns.style.display = 'none';
    }

    await loadTickets();

    const projectionId = getProjectionIdFromURL();
    if (projectionId) {
        await loadProjection(projectionId);
    } else {
        renderSeatMap(); // default 8×10
    }
});