import React, { ReactNode, useState } from "react";
import { updateData } from "./updateData";
import { Dict } from "../../api/api";

type SavingStatus = "untouched" | "touched" | "saving";

export interface FormProps {
  label: string;
  children: ReactNode;
  onDelete: () => Promise<void>;
  onSave: () => Promise<void>;
  saving: SavingStatus;
  onUp: () => Promise<void>;
  onDown: () => Promise<void>;
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
  onDelete: () => Promise<void>;
  onUp: () => Promise<void>;
  onDown: () => Promise<void>;
}

export const DynamicForm: React.FC<DynamicFormProps> = ({
  name,
  data,
  setData,
  ui,
  onSave,
  onDelete,
  onUp,
  onDown,
}) => {
  const [saving, setSaving] = useState<SavingStatus>("untouched");
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
              setSaving("touched");
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
            setSaving("touched");
            setData(updateData(data, path, newValue));
          }}
        />
      );
    });
  }

  return (
    <ui.Form
      label={name}
      saving={saving}
      onSave={async () => {
        setSaving("saving");
        await onSave();
        setSaving("untouched");
      }}
      onDelete={onDelete}
      onUp={onUp}
      onDown={onDown}
    >
      {renderNode(data)}
    </ui.Form>
  );
};
