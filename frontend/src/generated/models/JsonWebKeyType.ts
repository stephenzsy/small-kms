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
export const JsonWebKeyType = {
    Ec: 'EC',
    Rsa: 'RSA',
    Oct: 'oct'
} as const;
export type JsonWebKeyType = typeof JsonWebKeyType[keyof typeof JsonWebKeyType];


export function JsonWebKeyTypeFromJSON(json: any): JsonWebKeyType {
    return JsonWebKeyTypeFromJSONTyped(json, false);
}

export function JsonWebKeyTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): JsonWebKeyType {
    return json as JsonWebKeyType;
}

export function JsonWebKeyTypeToJSON(value?: JsonWebKeyType | null): any {
    return value as any;
}

