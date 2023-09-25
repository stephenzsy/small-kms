/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
import type { NamespaceProfile } from './NamespaceProfile';
import {
    NamespaceProfileFromJSON,
    NamespaceProfileFromJSONTyped,
    NamespaceProfileToJSON,
} from './NamespaceProfile';

/**
 * 
 * @export
 * @interface MyProfile
 */
export interface MyProfile {
    /**
     * 
     * @type {NamespaceProfile}
     * @memberof MyProfile
     */
    user?: NamespaceProfile;
    /**
     * 
     * @type {NamespaceProfile}
     * @memberof MyProfile
     */
    device?: NamespaceProfile;
}

/**
 * Check if a given object implements the MyProfile interface.
 */
export function instanceOfMyProfile(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function MyProfileFromJSON(json: any): MyProfile {
    return MyProfileFromJSONTyped(json, false);
}

export function MyProfileFromJSONTyped(json: any, ignoreDiscriminator: boolean): MyProfile {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'user': !exists(json, 'user') ? undefined : NamespaceProfileFromJSON(json['user']),
        'device': !exists(json, 'device') ? undefined : NamespaceProfileFromJSON(json['device']),
    };
}

export function MyProfileToJSON(value?: MyProfile | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'user': NamespaceProfileToJSON(value.user),
        'device': NamespaceProfileToJSON(value.device),
    };
}
