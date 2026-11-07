
window.initKanbanBoard = function(boardElementOrSelector, options = {}) {
    const board = typeof boardElementOrSelector === 'string'
        ? document.querySelector(boardElementOrSelector)
        : boardElementOrSelector;

    if (!board) return [];

    const groupName = options.group || board.getAttribute('data-kanban-group') || 'kanban';
    const postUrl = options.postUrl || board.getAttribute('data-kanban-url') || '';
    const columnSelector = options.columnSelector || '.crm-kanban-column-body';

    const columns = board.querySelectorAll(columnSelector);
    const instances = [];

    columns.forEach(col => {
        const sortable = new Sortable(col, {
            group: groupName,
            animation: 150,
            ghostClass: 'crm-kanban-ghost',
            chosenClass: 'crm-kanban-chosen',
            dragClass: 'crm-kanban-dragging',
            handle: options.handle || undefined,
            onEnd: function(evt) {
                const itemEl = evt.item;
                const toColumn = evt.to;
                const fromColumn = evt.from;
                const itemId = itemEl.getAttribute('data-id') || itemEl.id;
                const toColId = toColumn.getAttribute('data-col-id') || toColumn.id.replace('kanban-items-', '');
                const fromColId = fromColumn.getAttribute('data-col-id') || fromColumn.id.replace('kanban-items-', '');
                const newIndex = evt.newIndex;

                updateColumnBadges(fromColumn, toColumn);

                const customEvent = new CustomEvent('kanban:cardMoved', {
                    bubbles: true,
                    detail: {
                        itemId,
                        toColId,
                        fromColId,
                        newIndex,
                        itemEl,
                        toColumn,
                        fromColumn
                    }
                });
                board.dispatchEvent(customEvent);

                const endpoint = postUrl || itemEl.getAttribute('hx-post') || toColumn.getAttribute('data-post-url');
                if (endpoint && window.htmx) {
                    window.htmx.ajax('POST', endpoint, {
                        values: {
                            id: itemId,
                            to_stage: toColId,
                            from_stage: fromColId,
                            position: newIndex
                        },
                        swap: 'none'
                    });
                }
            }
        });
        instances.push(sortable);
    });

    function updateColumnBadges(fromCol, toCol) {
        [fromCol, toCol].forEach(col => {
            if (!col) return;
            const parentCol = col.closest('.crm-kanban-column');
            if (!parentCol) return;
            const badge = parentCol.querySelector('.crm-kanban-column-count');
            if (badge) {
                const count = col.querySelectorAll('.crm-kanban-card').length;
                badge.textContent = count;
            }
        });
    }

    return instances;
};

document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('.crm-kanban-board').forEach(board => {
        window.initKanbanBoard(board);
    });
});

