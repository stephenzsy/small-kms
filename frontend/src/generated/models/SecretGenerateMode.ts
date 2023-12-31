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
export const SecretGenerateMode = {
    SecretGenerateModeManual: 'manual',
    SecretGenerateModeServerGeneratedRandom: 'random-server'
} as const;
export type SecretGenerateMode = typeof SecretGenerateMode[keyof typeof SecretGenerateMode];


export function SecretGenerateModeFromJSON(json: any): SecretGenerateMode {
    return SecretGenerateModeFromJSONTyped(json, false);
}

export function SecretGenerateModeFromJSONTyped(json: any, ignoreDiscriminator: boolean): SecretGenerateMode {
    return json as SecretGenerateMode;
}

export function SecretGenerateModeToJSON(value?: SecretGenerateMode | null): any {
    return value as any;
}

