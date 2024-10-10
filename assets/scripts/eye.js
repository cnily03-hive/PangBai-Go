!function () {
    const searchParams = new URLSearchParams(window.location.search);
    let editIndex = searchParams.get('edit');
    if (searchParams.has('edit')) {
        searchParams.delete('edit');
        window.history.replaceState(null, null, window.location.pathname + (searchParams.size > 0 ? '?' + searchParams.toString() : ''));
    }
    document.addEventListener('DOMContentLoaded', function () {
        const boxes = document.querySelectorAll('.box');
        for (let i = 0; i < 2; ++i) {
            boxes[i].addEventListener('mouseover', function () {
                boxes[i].classList.add('hover');
                boxes[1 - i].classList.remove('hover');
            });
        }
        let onkeydown = false
        const contenteditable = document.querySelector('[contenteditable]');
        if (!/[^\d]/.test(editIndex)) {
            editIndex = parseInt(editIndex);
            contenteditable.focus();
            const range = document.createRange();
            const sel = window.getSelection();
            const text = contenteditable.innerText;
            editIndex = Math.min(editIndex, text.length);
            editIndex = Math.max(editIndex, 0);
            if (!searchParams.get('input')) editIndex = text.length;
            range.setStart(contenteditable.childNodes[0], editIndex);
            range.collapse(true);
            sel.removeAllRanges();
            sel.addRange(range);
        }
        contenteditable.addEventListener('keydown', function (e) {
            if (onkeydown) { return; }
            onkeydown = true
            if (e.keyCode === 13) {
                e.preventDefault();
                searchParams.set('input', contenteditable.innerText);
                const cursorIndex = window.getSelection().anchorOffset;
                searchParams.set('edit', cursorIndex);
                window.location.search = searchParams.toString();
            }
        });
        contenteditable.addEventListener('keyup', function (e) {
            onkeydown = false
        });
    })
}()