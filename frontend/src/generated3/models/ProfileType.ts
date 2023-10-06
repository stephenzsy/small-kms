/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API Shared
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
export const ProfileType = {
    ProfileTypeRootCA: 'root-ca',
    ProfileTypeIntermediateCA: 'intermediate-ca',
    ProfileTypeServicePrincipal: 'service-principal',
    ProfileTypeGroup: 'group',
    ProfileTypeDevice: 'device',
    ProfileTypeUser: 'user',
    ProfileTypeApplication: 'application'
} as const;
export type ProfileType = typeof ProfileType[keyof typeof ProfileType];


export function ProfileTypeFromJSON(json: any): ProfileType {
    return ProfileTypeFromJSONTyped(json, false);
}

export function ProfileTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProfileType {
    return json as ProfileType;
}

export function ProfileTypeToJSON(value?: ProfileType | null): any {
    return value as any;
}

