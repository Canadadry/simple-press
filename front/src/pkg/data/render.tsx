import React, { ReactNode } from "react";
import { updateData } from "./updateData";
import { Dict } from "../../api/api";

export interface FormProps {
  label: string;
  children: ReactNode;
  onSave: () => Promise<void>;
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
  setData: (fullPath: string, value: string) => void;
}

export interface FormCheckBoxProps {
  label: string;
  name: string;
  checked: boolean;
  setData: (fullPath: string, value: string) => void;
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
  setData: (data: Dict) => void;
  ui: DynamicFormUI;
  onSave: () => Promise<void>;
}

export const DynamicForm: React.FC<DynamicFormProps> = ({
  name,
  data,
  setData,
  ui,
  onSave,
}) => {
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

      if (typeof value === "boolean") {
        return (
          <ui.FormCheckBox
            key={fullPath}
            label={key}
            name={fullPath}
            checked={value}
            setData={(path, newValue) => {
              setData(updateData(data, path, newValue));
            }}
          />
        );
      }

      return (
        <ui.FormInput
          key={fullPath}
          label={key}
          name={fullPath}
          inputType={typeof value === "number" ? "number" : "text"}
          value={String(value)}
          setData={(path, newValue) => {
            setData(updateData(data, path, newValue));
          }}
        />
      );
    });
  }

  return (
    <ui.Form label={name} onSave={onSave}>
      {renderNode(data)}
    </ui.Form>
  );
};
