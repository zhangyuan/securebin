const hexToUint8Array = (hex: string) => {
  const bytes = [];
  for (let i = 0; i < hex.length; i += 2) {
    bytes.push(parseInt(hex.substring(i, i + 2), 16));
  }
  return new Uint8Array(bytes);
};

export const secureRandomString = (length: number): string => {
  const array = new Uint8Array(Math.ceil(length / 2));
  window.crypto.getRandomValues(array);
  let result = '';
  for (let i = 0; i < array.length; i++) {
      result += ('0' + array[i].toString(16)).slice(-2);
  }
  return result.slice(0, length);
}

const arrayBufferToHex = (arrayBuffer: ArrayBuffer): string => {
  const uint8Array = new Uint8Array(arrayBuffer);
  return Array.from(uint8Array)
     .map((byte) => byte.toString(16).padStart(2, '0'))
     .join('');
}

export const encrypt = async (
  plainText: string,
  keyHex: string,
  ivHex: string,
): Promise<string> => {
  const encoder = new TextEncoder();
  const data = encoder.encode(plainText);
  const key = hexToUint8Array(keyHex);
  const iv = hexToUint8Array(ivHex);
  const cryptoKey = await window.crypto.subtle.importKey(
    "raw",
    key,
    { name: "AES-GCM", length: 256 },
    false,
    ["encrypt"],
  );
  const cipherText = await window.crypto.subtle.encrypt(
    {
      name: "AES-GCM",
      iv: iv,
    },
    cryptoKey,
    data,
  );
  return arrayBufferToHex(cipherText);
};

export const decrypt = async (
  cipherTextHex: string,
  keyHex: string,
  ivHex: string,
) => {
  const key = hexToUint8Array(keyHex);
  const iv = hexToUint8Array(ivHex);
  const cipherText = hexToUint8Array(cipherTextHex);

  const cryptoKey = await window.crypto.subtle.importKey(
    "raw",
    key,
    { name: "AES-GCM", length: 256 },
    false,
    ["decrypt"],
  );
  const plainBuffer = await window.crypto.subtle.decrypt(
    {
      name: "AES-GCM",
      iv: iv,
    },
    cryptoKey,
    cipherText,
  );
  const decoder = new TextDecoder();
  const plainText = decoder.decode(plainBuffer);
  return plainText;
};
