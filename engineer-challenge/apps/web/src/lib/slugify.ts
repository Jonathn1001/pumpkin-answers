// slugify turns a tenant display name into a URL-safe slug that satisfies the
// backend slug rule ^[a-z0-9][a-z0-9-]{0,62}$ (see internal/usecase/tenant.go).
// Diacritics (incl. Vietnamese) are stripped so "Bảo Việt" → "bao-viet".
export function slugify(name: string): string {
  return name
    .normalize("NFD")
    .replace(/[̀-ͯ]/g, "") // drop combining diacritical marks
    .replace(/[đĐ]/g, "d") // đ/Đ don't decompose under NFD
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-") // non-alphanumeric runs → single hyphen
    .replace(/^-+|-+$/g, "") // trim leading/trailing hyphens
    .slice(0, 63) // cap to backend max length
    .replace(/-+$/g, ""); // re-trim if slice landed mid-hyphen
}
