export function toPEMBlock(encoded: string, type: string) {
  const lines = encoded.match(/.{1,64}/g) ?? [];
  return `-----BEGIN ${type}-----\n${lines.join("\n")}\n-----END ${type}-----`;
}

export function base64UrlEncodedToStdEncoded(data: string) {
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

export function base64StdEncodedToUrlEncoded(data: string) {
  return data.replace(/\+/g, "-").replace(/\//g, "_").replace(/=/g, "");
}

export function base64UrlEncodeBuffer(buffer: ArrayBuffer): string {
  return base64StdEncodedToUrlEncoded(
    btoa(String.fromCharCode(...new Uint8Array(buffer)))
  );
}

export function base64UrlDecodeBuffer(text: string): ArrayBuffer {
  const s = atob(base64UrlEncodedToStdEncoded(text));
  const bytes = new Uint8Array(s.length);
  for (let i = 0; i < s.length; i++) {
    bytes[i] = s.charCodeAt(i);
  }
  return bytes.buffer;
}
