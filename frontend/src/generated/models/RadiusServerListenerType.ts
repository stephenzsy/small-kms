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
export const RadiusServerListenerType = {
    RadiusServerListenerTypeAuth: 'auth',
    RadiusServerListenerTypeAcct: 'acct'
} as const;
export type RadiusServerListenerType = typeof RadiusServerListenerType[keyof typeof RadiusServerListenerType];


export function RadiusServerListenerTypeFromJSON(json: any): RadiusServerListenerType {
    return RadiusServerListenerTypeFromJSONTyped(json, false);
}

export function RadiusServerListenerTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): RadiusServerListenerType {
    return json as RadiusServerListenerType;
}

export function RadiusServerListenerTypeToJSON(value?: RadiusServerListenerType | null): any {
    return value as any;
}

