export function formatBytes(bytes: number, decimals = 2): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

// Optimize date formatting performance by reusing the formatter instance.
// This avoids creating new Intl.DateTimeFormat objects and parsing locale data on every render,
// significantly improving performance for long lists of dates.
const dateFormatter = new Intl.DateTimeFormat(undefined, {
  year: 'numeric',
  month: 'numeric',
  day: 'numeric',
});

/**
 * Optimized date formatting using cached Intl.DateTimeFormat instance.
 * Much faster than new Date().toLocaleDateString() in loops.
 * Handles string, number, or Date objects. Returns 'Invalid Date' if parsing fails.
 */
export function formatDate(date: string | number | Date): string {
  if (!date) return '';
  const d = date instanceof Date ? date : new Date(date);
  if (isNaN(d.getTime())) return 'Invalid Date';
  return dateFormatter.format(d);
}
