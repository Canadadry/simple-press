import React, { ReactNode, useEffect, useState, useCallback } from "react";
import { updateData } from "./updateData";
import { Dict } from "../../api/api";
import { TextArea } from "@radix-ui/themes";

type SavingStatus = "untouched" | "touched" | "saving";
type Mode = "json" | "form";

export interface FormProps {
  label: string;
  children: ReactNode;
  mode: Mode;
  setMode: (m: Mode) => void;
  onSave: () => Promise<void>;
  saving: SavingStatus;
  onDelete?: () => Promise<void>;
  onUp?: () => Promise<void>;
  onDown?: () => Promise<void>;
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
  onSave: (data: Dict) => Promise<void>;
  onDelete?: () => Promise<void>;
  onUp?: () => Promise<void>;
  onDown?: () => Promise<void>;
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
  const [mode, setMode] = useState<Mode>("form");
  const [temp, setTemp] = useState<string | null>(null);
  const [cache, setCache] = useState<{ key: string; value: string }[]>([]);
  useEffect(() => {
    if (mode === "json") {
      cache.forEach((c: { key: string; value: string }) => {
        data = updateData(data, c.key, c.value);
      });
      setCache([]);
      setTemp(JSON.stringify(data, null, 2));
      return;
    }
    if (temp !== null) {
      try {
        const parsed = JSON.parse(temp);
        setData(parsed);
        setSaving("touched");
      } catch {
        // JSON invalide, on ignore
      }
    }
  }, [mode]);

  const renderNode = useCallback(
    (obj: Dict, prefix: string = ""): React.ReactNode => {
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
                setCache(cache.concat({ key: path, value: newValue }));
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
              setCache(cache.concat({ key: path, value: newValue }));
            }}
          />
        );
      });
    },
    [setSaving, cache, setCache, ui],
  );

  return (
    <ui.Form
      mode={mode}
      setMode={setMode}
      label={name}
      saving={saving}
      onSave={async () => {
        setSaving("saving");
        cache.forEach((c: { key: string; value: string }) => {
          data = updateData(data, c.key, c.value);
        });
        setCache([]);
        await onSave(data);
        setSaving("untouched");
      }}
      onDelete={onDelete}
      onUp={onUp}
      onDown={onDown}
    >
      {mode === "form" ? (
        renderNode(data)
      ) : (
        <TextArea
          spellCheck={false}
          variant="soft"
          rows={10}
          value={temp || ""}
          disabled={saving === "saving"}
          onChange={(e) => {
            setTemp(e.target.value);
            let p: Dict | null = null;
            try {
              p = JSON.parse(e.target.value);
            } catch {
              setSaving("untouched");
            } finally {
              if (p != null) {
                // setData(p);
                setSaving("touched");
              }
            }
          }}
        />
      )}
    </ui.Form>
  );
};
