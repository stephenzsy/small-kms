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
 * @interface ResourceLocator2
 */
export interface ResourceLocator2 {
    /**
     * 
     * @type {NamespaceKind}
     * @memberof ResourceLocator2
     */
    namespaceKind: NamespaceKind;
    /**
     * 
     * @type {string}
     * @memberof ResourceLocator2
     */
    namespaceIdentifier: string;
    /**
     * 
     * @type {ResourceKind}
     * @memberof ResourceLocator2
     */
    resourceKind: ResourceKind;
    /**
     * 
     * @type {string}
     * @memberof ResourceLocator2
     */
    resourceIdentifier: string;
}

/**
 * Check if a given object implements the ResourceLocator2 interface.
 */
export function instanceOfResourceLocator2(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "namespaceKind" in value;
    isInstance = isInstance && "namespaceIdentifier" in value;
    isInstance = isInstance && "resourceKind" in value;
    isInstance = isInstance && "resourceIdentifier" in value;

    return isInstance;
}

export function ResourceLocator2FromJSON(json: any): ResourceLocator2 {
    return ResourceLocator2FromJSONTyped(json, false);
}

export function ResourceLocator2FromJSONTyped(json: any, ignoreDiscriminator: boolean): ResourceLocator2 {
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

export function ResourceLocator2ToJSON(value?: ResourceLocator2 | null): any {
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

