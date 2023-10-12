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

import { exists, mapValues } from '../runtime';
import type { NamespaceKind } from './NamespaceKind';
import {
    NamespaceKindFromJSON,
    NamespaceKindFromJSONTyped,
    NamespaceKindToJSON,
} from './NamespaceKind';

/**
 * 
 * @export
 * @interface ProfileRefFields
 */
export interface ProfileRefFields {
    /**
     * 
     * @type {NamespaceKind}
     * @memberof ProfileRefFields
     */
    type: NamespaceKind;
    /**
     * Display name of the resource
     * @type {string}
     * @memberof ProfileRefFields
     */
    displayName: string;
    /**
     * Whether the resource is managed by the application
     * @type {boolean}
     * @memberof ProfileRefFields
     */
    isAppManaged?: boolean;
}

/**
 * Check if a given object implements the ProfileRefFields interface.
 */
export function instanceOfProfileRefFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "type" in value;
    isInstance = isInstance && "displayName" in value;

    return isInstance;
}

export function ProfileRefFieldsFromJSON(json: any): ProfileRefFields {
    return ProfileRefFieldsFromJSONTyped(json, false);
}

export function ProfileRefFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProfileRefFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'type': NamespaceKindFromJSON(json['type']),
        'displayName': json['displayName'],
        'isAppManaged': !exists(json, 'isAppManaged') ? undefined : json['isAppManaged'],
    };
}

export function ProfileRefFieldsToJSON(value?: ProfileRefFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'type': NamespaceKindToJSON(value.type),
        'displayName': value.displayName,
        'isAppManaged': value.isAppManaged,
    };
}

