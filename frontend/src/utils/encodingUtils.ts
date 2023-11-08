
export function toPEMBlock(encoded: string, type: string) {
  const lines = encoded.match(/.{1,64}/g) ?? [];
  return `-----BEGIN ${type}-----\n${lines.join("\n")}\n-----END ${type}-----`;
}export function base64UrlEncodedToStdEncoded(data: string) {
  const r = data.replace(/-/g, "+").replace(/_/g, "/");
  // decide pad length
  const padLen = r.length % 4;
  if (padLen == 2) {
    return r + "==";
  } else if (padLen == 3) {
    return r + "=";
  }
  return r;
}

