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

import { exists, mapValues } from '../runtime';
import type { ProfileType } from './ProfileType';
import {
    ProfileTypeFromJSON,
    ProfileTypeFromJSONTyped,
    ProfileTypeToJSON,
} from './ProfileType';
import type { ResourceMetadata } from './ResourceMetadata';
import {
    ResourceMetadataFromJSON,
    ResourceMetadataFromJSONTyped,
    ResourceMetadataToJSON,
} from './ResourceMetadata';

/**
 * 
 * @export
 * @interface ProfileRef
 */
export interface ProfileRef {
    /**
     * 
     * @type {ProfileType}
     * @memberof ProfileRef
     */
    type: ProfileType;
    /**
     * Identifier of the resource
     * @type {string}
     * @memberof ProfileRef
     */
    id: string;
    /**
     * Display name of the resource
     * @type {string}
     * @memberof ProfileRef
     */
    displayName: string;
    /**
     * 
     * @type {ResourceMetadata}
     * @memberof ProfileRef
     */
    metadata?: ResourceMetadata;
}

/**
 * Check if a given object implements the ProfileRef interface.
 */
export function instanceOfProfileRef(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "type" in value;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "displayName" in value;

    return isInstance;
}

export function ProfileRefFromJSON(json: any): ProfileRef {
    return ProfileRefFromJSONTyped(json, false);
}

export function ProfileRefFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProfileRef {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'type': ProfileTypeFromJSON(json['type']),
        'id': json['id'],
        'displayName': json['displayName'],
        'metadata': !exists(json, 'metadata') ? undefined : ResourceMetadataFromJSON(json['metadata']),
    };
}

export function ProfileRefToJSON(value?: ProfileRef | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'type': ProfileTypeToJSON(value.type),
        'id': value.id,
        'displayName': value.displayName,
        'metadata': ResourceMetadataToJSON(value.metadata),
    };
}
