/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.1
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


/**
 * 
 * @export
 */
export const SecretRandomCharacterClass = {
    SecretRandomCharClassBase64RawURL: 'base64-raw-url'
} as const;
export type SecretRandomCharacterClass = typeof SecretRandomCharacterClass[keyof typeof SecretRandomCharacterClass];


export function SecretRandomCharacterClassFromJSON(json: any): SecretRandomCharacterClass {
    return SecretRandomCharacterClassFromJSONTyped(json, false);
}

export function SecretRandomCharacterClassFromJSONTyped(json: any, ignoreDiscriminator: boolean): SecretRandomCharacterClass {
    return json as SecretRandomCharacterClass;
}

export function SecretRandomCharacterClassToJSON(value?: SecretRandomCharacterClass | null): any {
    return value as any;
}

