// Reveal the current page inside the scrollable tree without moving the page.
const sectionNavigation = document.querySelector("#td-section-nav");

function revealActiveNavigation() {
    const active = sectionNavigation?.querySelector(".td-sidebar-link.active");
    if (!active || !sectionNavigation.clientHeight) return;
    const row = active.getBoundingClientRect();
    const tree = sectionNavigation.getBoundingClientRect();
    if (row.top < tree.top || row.bottom > tree.bottom) {
        sectionNavigation.scrollTop +=
            row.top -
            tree.top -
            (sectionNavigation.clientHeight - row.height) / 2;
    }
}

revealActiveNavigation();
// Docsy hydrates cached trees at DOMContentLoaded before this listener runs.
document.addEventListener("DOMContentLoaded", revealActiveNavigation, {
    once: true,
});
sectionNavigation?.addEventListener(
    "shown.bs.collapse",
    revealActiveNavigation,
);
