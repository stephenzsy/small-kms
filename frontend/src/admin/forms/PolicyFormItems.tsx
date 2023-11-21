import { Checkbox, Form, FormInstance, FormItemProps, Radio } from "antd";
import {
  JsonWebKeyCurveName,
  JsonWebKeyOperation,
  JsonWebKeyType,
} from "../../generated/apiv2";
import { useWatch } from "antd/es/form/Form";
import { CheckboxChangeEvent } from "antd/es/checkbox";

export type KeySpecFormItemsProps<T> = {
  formInstance: FormInstance<T>;
  ktyName: FormItemProps<T>["name"];
  keySizeName: FormItemProps<T>["name"];
  crvName: FormItemProps<T>["name"];
  keyOpsName: FormItemProps<T>["name"];
};

export function KeySpecFormItems<T>({
  ktyName,
  keySizeName,
  crvName,
  keyOpsName,
  formInstance,
}: KeySpecFormItemsProps<T>) {
  const kty = useWatch(ktyName, formInstance);
  return (
    <>
      <Form.Item label="Key Type" name={ktyName}>
        <Radio.Group>
          <Radio value={JsonWebKeyType.Rsa}>RSA</Radio>
          <Radio value={JsonWebKeyType.Ec}>EC</Radio>
        </Radio.Group>
      </Form.Item>
      {kty === JsonWebKeyType.Rsa && (
        <Form.Item<T> label="Key size" name={keySizeName}>
          <Radio.Group>
            <Radio value={2048}>2048</Radio>
            <Radio value={3072}>3072</Radio>
            <Radio value={4096}>4096</Radio>
          </Radio.Group>
        </Form.Item>
      )}
      {kty === JsonWebKeyType.Ec && (
        <Form.Item<T> label="Curve name" name={crvName}>
          <Radio.Group>
            <Radio value={JsonWebKeyCurveName.CurveNameP256}>P-256</Radio>
            <Radio value={JsonWebKeyCurveName.CurveNameP384}>P-384</Radio>
            <Radio value={JsonWebKeyCurveName.CurveNameP521}>P-521</Radio>
          </Radio.Group>
        </Form.Item>
      )}
      <Form.Item<T> name={keyOpsName} label="Key Operations">
        <Checkbox.Group className="inline-grid grid-cols-[auto_auto] gap-4">
          <Checkbox value={JsonWebKeyOperation.Sign}>Sign</Checkbox>
          <Checkbox value={JsonWebKeyOperation.Verify}>Verify</Checkbox>
          <Checkbox value={JsonWebKeyOperation.Encrypt}>Encrypt</Checkbox>
          <Checkbox value={JsonWebKeyOperation.Decrypt}>Decrypt</Checkbox>
          <Checkbox value={JsonWebKeyOperation.WrapKey}>Wrap Key</Checkbox>
          <Checkbox value={JsonWebKeyOperation.UnwrapKey}>Unwrap Key</Checkbox>
          <Checkbox value={JsonWebKeyOperation.DeriveKey}>Derive Key</Checkbox>
          <Checkbox value={JsonWebKeyOperation.DeriveBits}>
            Derive Bits
          </Checkbox>
        </Checkbox.Group>
      </Form.Item>
    </>
  );
}

export function KeyExportableFormItem<T>({
  name,
  disabled,
}: {
  name: FormItemProps<T>["name"];
  disabled?: boolean;
}) {
  return (
    <Form.Item<T>
      name={name}
      valuePropName="checked"
      getValueFromEvent={(e: CheckboxChangeEvent) => {
        return e.target.checked;
      }}
    >
      <Checkbox disabled={disabled}>Private Key Exportable</Checkbox>
    </Form.Item>
  );
}
