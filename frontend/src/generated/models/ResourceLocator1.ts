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
import type { ResourceKind } from './ResourceKind';
import {
    ResourceKindFromJSON,
    ResourceKindFromJSONTyped,
    ResourceKindToJSON,
} from './ResourceKind';

/**
 * 
 * @export
 * @interface ResourceLocator1
 */
export interface ResourceLocator1 {
    /**
     * 
     * @type {NamespaceKind}
     * @memberof ResourceLocator1
     */
    namespaceKind: NamespaceKind;
    /**
     * 
     * @type {string}
     * @memberof ResourceLocator1
     */
    namespaceIdentifier: string;
    /**
     * 
     * @type {ResourceKind}
     * @memberof ResourceLocator1
     */
    resourceKind: ResourceKind;
    /**
     * 
     * @type {string}
     * @memberof ResourceLocator1
     */
    resourceIdentifier: string;
}

/**
 * Check if a given object implements the ResourceLocator1 interface.
 */
export function instanceOfResourceLocator1(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "namespaceKind" in value;
    isInstance = isInstance && "namespaceIdentifier" in value;
    isInstance = isInstance && "resourceKind" in value;
    isInstance = isInstance && "resourceIdentifier" in value;

    return isInstance;
}

export function ResourceLocator1FromJSON(json: any): ResourceLocator1 {
    return ResourceLocator1FromJSONTyped(json, false);
}

export function ResourceLocator1FromJSONTyped(json: any, ignoreDiscriminator: boolean): ResourceLocator1 {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'namespaceKind': NamespaceKindFromJSON(json['namespaceKind']),
        'namespaceIdentifier': json['namespaceIdentifier'],
        'resourceKind': ResourceKindFromJSON(json['resourceKind']),
        'resourceIdentifier': json['resourceIdentifier'],
    };
}

export function ResourceLocator1ToJSON(value?: ResourceLocator1 | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'namespaceKind': NamespaceKindToJSON(value.namespaceKind),
        'namespaceIdentifier': value.namespaceIdentifier,
        'resourceKind': ResourceKindToJSON(value.resourceKind),
        'resourceIdentifier': value.resourceIdentifier,
    };
}
