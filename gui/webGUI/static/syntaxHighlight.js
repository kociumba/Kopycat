const LogHighlighter = (() => {
    // Compile regex pattern once
    const pattern = /^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})(\s*)(INFO|WARN|ERRO|DEBU)(\s*)(<[^>]+>)(\s*)(.+?)(\s*)$/;

    // Function to escape HTML special characters
    const escapeHtml = (unsafe) => {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    };

    return {
        highlightLog: function(logContent) {
            return logContent.split('\n').map(line => {
                const match = pattern.exec(line);
                if (match) {
                    const [_, datetime, ws1, level, ws2, module, ws3, message, ws4] = match;
                    const [date, time] = datetime.split(' ');
                    // return `<span class="module">${lineNumber}  </span>` +
                    return `<span class="date">${date} ${time}${ws1}</span>` +
                           `<span class="level-${level.toLowerCase()}">${level}${ws2}</span>` +
                           `<span class="module">${escapeHtml(module)}${ws3}</span>` +
                           `<span class="message">${escapeHtml(message)}${ws4}</span>`;
                }
                return escapeHtml(line); // Return escaped line if it doesn't match the pattern
            }).join('\n');
        }
    };
})();
