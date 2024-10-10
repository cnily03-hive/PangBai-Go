

document.addEventListener('DOMContentLoaded', function () {
    const material = Object.freeze({
        0: 'PangBai 出院有一些日子了，加上住院期间经常看 <strong>MyGO!!!!!</strong> 的缘故，这些天来他心情比较低落，但今天已经能正常打 <strong>GoGo</strong> 了，看起来和她所说的一样，她是「高性能」的呢！',
        1: '可是当我买菜回来时，却发现她倒在地上，眉毛的样子像两个<strong>大括号</strong>，左眼呈现出「.User」，而右眼则是「PangBai」。'
    })

    const selectbox = document.querySelector('.select-box');
    let cont = document.querySelector('.box');
    let textspan = document.getElementById('text');
    const attachEvent = () => cont.classList.add('event')
    const detachEvent = () => cont.classList.remove('event')
    const setAutoWidth = () => {
        cont.classList.add('justify');
        textspan.style.width = ''
        textspan.style.height = ''
    }
    const setFitWidth = (html) => {
        textspan.style.opacity = '0';
        setAutoWidth();
        cont.classList.add('justify');
        textspan.innerHTML = html;
        let width = textspan.offsetWidth;
        let height = textspan.offsetHeight;
        textspan.style.width = width + 'px';
        textspan.style.height = height + 'px';
        textspan.innerHTML = '';
        textspan.style.opacity = '1';
    }
    const setText = (html) => {
        textspan.innerHTML = html;
    }
    const getHtmlText = (html) => {
        const sp = document.createElement('span');
        sp.innerHTML = html;
        return sp.innerText;
    }
    let on_animate = false;
    let interval_array = [];
    const clearAllInterval = () => { while (interval_array.length) clearInterval(interval_array.pop()); }
    const animateText = (html, cb = () => { }) => {
        on_animate = true;
        clearAllInterval();
        setFitWidth(html);
        let appended = '';
        let pos = 0;
        let cnt = 0;
        let fulltext = getHtmlText(html);
        let interval = setInterval(() => {
            let c = html[pos];
            if (c === '<') {
                let end = pos;
                while (html[end] !== '>') end++;
                c = html.substring(pos, end + 1);
                pos = end;
            }
            appended += c;
            setText(appended);
            pos++;
            cnt++;
            if (cnt > 0 && cnt >= fulltext.length * 0.75) {
                on_animate = false;
            }
            if (pos >= html.length) {
                clearAllInterval();
                on_animate = false;
                if (typeof cb === 'function') cb();
            }
        }, 50)
        interval_array.push(interval);
    }
    const after_animate = () => {
        setTimeout(() => {
            selectbox.classList.replace('prevent', 'active');
        }, 1000)
    }
    let next_material = 0, TOT_MATERIAL = Object.keys(material).length;

    cont.addEventListener('click', function () {
        if (next_material >= TOT_MATERIAL) return;
        if (on_animate) return;
        animateText(material[next_material], next_material === TOT_MATERIAL - 1 ? after_animate : null);
        next_material++;
        if (next_material >= TOT_MATERIAL) {
            detachEvent();
        }
    })
})