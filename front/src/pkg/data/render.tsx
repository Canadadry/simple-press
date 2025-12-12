import React, { ReactNode } from "react";
import { Dict } from "./parseFormData";

export interface FormProps {
  label: string;
  children: ReactNode;
}

export interface FormObjectProps {
  label: string;
  children: ReactNode;
}

export interface FormInputProps {
  label: string;
  name: string;
  inputType: string;
  value: string;
}

export interface FormCheckBoxProps {
  label: string;
  name: string;
  checked: boolean;
}

export interface DynamicFormUI {
  Form: React.FC<FormProps>;
  FormObject: React.FC<FormObjectProps>;
  FormInput: React.FC<FormInputProps>;
  FormCheckBox: React.FC<FormCheckBoxProps>;
}

export interface DynamicFormProps {
  name: string;
  data: Dict;
  ui: DynamicFormUI;
}

export const DynamicForm: React.FC<DynamicFormProps> = ({ name, data, ui }) => {
  function renderNode(obj: Dict, prefix: string = ""): React.ReactNode {
    return Object.entries(obj).map(([key, value]) => {
      const fullPath = prefix ? `${prefix}.${key}` : key;

      if (Array.isArray(value)) {
        throw new Error(`Arrays not supported at path ${fullPath}`);
      }

      if (typeof value === "object" && value !== null) {
        return (
          <ui.FormObject key={fullPath} label={key}>
            {renderNode(value as Dict, fullPath)}
          </ui.FormObject>
        );
      }

      const inputType = typeof value === "number" ? "number" : "text";

      return (
        <ui.FormInput
          key={fullPath}
          label={key}
          name={fullPath}
          inputType={inputType}
          value={String(value)}
        />
      );
    });
  }

  return <ui.Form label={name}>{renderNode(data)}</ui.Form>;
};
