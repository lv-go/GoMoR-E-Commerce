import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";
import { useEffect, useMemo, useState } from "react";
import { useFieldArray, useForm, type RegisterOptions } from "react-hook-form";
import { FaChevronLeft, FaChevronRight, FaEdit, FaPlus, FaSave, FaTimes, FaTrash } from "react-icons/fa";
import { useRevalidator, useSearchParams } from "react-router";
import type { Page, PageRequest } from "~/schemas/api";

export type Column<T> = {
  key: string;
  header: string;
  content: (item: T) => React.ReactNode;
  editor?: (register: any) => React.ReactNode;
  registerOptions?: RegisterOptions<{ items: any[] }, any>;
};

export default function Table<T, ID>({
  header,
  newItem,
  getItemId,
  columns,
  useGetPage,
  useCreate,
  useUpdate,
  useDeleteById,
}: {
  header: string;
  newItem: () => T;
  getItemId: (item: T) => ID;
  columns: Column<T>[];
  useGetPage: (pageRequest: PageRequest) => UseQueryResult<Page<any>, Error>;
  useCreate: () => UseMutationResult<T, Error, Partial<T>>;
  useUpdate: () => UseMutationResult<T, Error, { id: ID; data: Partial<T> }>;
  useDeleteById: () => UseMutationResult<void, Error, ID>;
}) {
  const [selectedItemsIds, setSelectedItemsIds] = useState<Set<any>>(new Set<any>());
  const [editingRows, setEditingRows] = useState<Set<number>>(new Set<number>());
  const [searchParams, setSearchParams] = useSearchParams();
  const [loading, setLoading] = useState(false);

  const pageRequest = useMemo(() => ({
    offset: parseInt(searchParams.get("offset") ?? "0"),
    limit: parseInt(searchParams.get("limit") ?? "10"),
  }), [searchParams]);

  const { data, isLoading, error, isFetched, isFetching } = useGetPage(pageRequest);
  const page = data ?? { items: [], total: 0, size: 10, page: 1, totalPages: 1 };

  useEffect(() => {
    console.log("isFetched", isFetched);
    console.log("isFetching", isFetching);
    console.log("page", page);
    if (isFetched) {
      reset({ items: page.items });
    }
  }, [isFetching]);

  const { control, reset, register, getValues } = useForm({
    defaultValues: { items: page.items },
  });

  const items = getValues("items") ?? [];

  const totalPages = useMemo(() => Math.ceil(page.total / page.size), [page.total, page.size]);

  const { mutateAsync: create } = useCreate();
  const { mutateAsync: update } = useUpdate();
  const { mutateAsync: deleteById } = useDeleteById();

  function setCurrentPage(currentPage: number) {
    searchParams.set("page", currentPage.toString());
    setSearchParams(searchParams);
    setSelectedItemsIds(new Set<any>());
    setEditingRows(new Set<number>());
  }

  function setPageSize(event: React.ChangeEvent<HTMLSelectElement>) {
    searchParams.set("size", event.target.value);
    searchParams.set("page", "1");
    setSearchParams(searchParams);
    setSelectedItemsIds(new Set<any>());
    setEditingRows(new Set<number>());
  }

  function handleAdd() {
    const item = newItem();
    reset({ items: [item, ...items] });
    handleStartEditing(0);
  }

  async function handleDeleteSelectedItems() {
    setLoading(true);
    await Promise.all(
      items
        .filter((item) => selectedItemsIds.has(getItemId(item)))
        .map((item) => deleteById(getItemId(item)))
    );

    reset({ items: items.filter((item) => !selectedItemsIds.has(getItemId(item))) });
    setSelectedItemsIds(new Set<any>());
    setLoading(false);
  }

  async function handleDeleteItem(id: ID) {
    setLoading(true);
    await deleteById(id);

    reset({ items: items.filter((item) => getItemId(item) !== id) });
    setSelectedItemsIds(new Set<any>());
    setLoading(false);
  }

  function handleStartEditing(index: number) {
    setEditingRows(new Set([...editingRows, index]));
  }

  async function handleSave(index: number) {
    setLoading(true);
    // const id = getItemId(items[index]);
    const item = getValues(`items.${index}`) as T;
    const id = getItemId(item);
    if (id) {
      await update({ id, data: item });
    } else {
      await create(item);
    }

    editingRows.delete(index);
    setEditingRows(new Set(editingRows));
    setLoading(false);
    reset();
  }

  function handleCancel(index: number) {
    // if adding new item, remove it
    if (!getItemId(items[index])) {
      reset({ items: items.filter((item) => getItemId(item) !== getItemId(items[index])) });
    }
    editingRows.delete(index);
    setEditingRows(new Set(editingRows));
  }

  function handleSelectAll(e: React.ChangeEvent<HTMLInputElement>) {
    const checked = e.target.checked;
    if (checked) {
      setSelectedItemsIds(new Set(items.map((item) => getItemId(item))));
    } else {
      setSelectedItemsIds(new Set<any>());
    }
  }

  const handleSelectItem = (id: any) => (e: React.ChangeEvent<HTMLInputElement>) => {
    const checked = e.target.checked;
    if (checked) {
      selectedItemsIds.add(id);
    } else {
      selectedItemsIds.delete(id);
    }
    setSelectedItemsIds(new Set(selectedItemsIds));
  }

  return <div className="flex flex-col h-full">
    <header className="navbar px-6 bg-gray-200 text-gray-800 border-b border-gray-300">
      <label htmlFor="root-drawer" aria-label="open sidebar" className="btn btn-square btn-ghost">
        {/* Sidebar toggle icon */}
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" strokeLinejoin="round" strokeLinecap="round" strokeWidth="2" fill="none" stroke="currentColor" className="my-1.5 inline-block size-4"><path d="M4 4m0 2a2 2 0 0 1 2 -2h12a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2z"></path><path d="M9 4v16"></path><path d="M14 10l2 2l-2 2"></path></svg>
      </label>
      <div className="flex items-center gap-2">
        <h2 className="text-xl font-semibold flex-1">{header}</h2>
      </div>
    </header>

    <main className="flex-1 overflow-auto">
      <table className="table">
        <thead>
          <tr>
            <th style={{ width: "56px" }}>
              <input type="checkbox" className="checkbox"
                aria-label="Select All" title="Select All"
                onChange={handleSelectAll} />
            </th>
            {columns.map((column) => (
              <th key={column.key}>{column.header}</th>
            ))}
            <th style={{ width: "112px" }}>
              <button className="btn btn-square"
                aria-label="Add" title="Add"
                disabled={items[0] && !getItemId(items[0])}
                onClick={handleAdd}>
                <FaPlus />
              </button>
              <button className="btn btn-square"
                aria-label="Delete Selected" title="Delete Selected"
                disabled={selectedItemsIds.size === 0}
                onClick={handleDeleteSelectedItems}>
                <FaTrash />
              </button>
            </th>
          </tr>
        </thead>
        <tbody>
          {items.map((item, index) => (
            <tr key={item.id}>
              <td>
                <input type="checkbox" className="checkbox"
                  aria-label="Select Category" title="Select Category"
                  value={item.id}
                  checked={selectedItemsIds.has(item.id ?? "")}
                  onChange={handleSelectItem(item.id)} />
              </td>
              {columns.map((column) => {
                return <td key={column.key}>{(editingRows.has(index))
                  ? column.editor?.(register(`items.${index}.${column.key}`, column.registerOptions))
                  : column.content(item)}</td>;
              }
              )}
              <td>
                {editingRows.has(index) ? (
                  <>
                    <button className="btn btn-square"
                      aria-label="Save" title="Save"
                      onClick={() => handleSave(index)}>
                      <FaSave />
                    </button>
                    <button className="btn btn-square"
                      aria-label="Cancel" title="Cancel"
                      onClick={() => handleCancel(index)}>
                      <FaTimes />
                    </button>
                  </>
                ) : (<>
                  <button className="btn btn-square"
                    aria-label="Edit" title="Edit"
                    onClick={() => handleStartEditing(index)}>
                    <FaEdit />
                  </button>
                  <button className="btn btn-square"
                    aria-label="Delete" title="Delete"
                    onClick={() => handleDeleteItem(getItemId(item))}>
                    <FaTrash />
                  </button>
                </>)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </main>

    <footer className="bg-gray-200 text-gray-800 p-2 border-t border-gray-300 flex ">
      <div className="flex-1 flex gap-2 items-center">
        <span className="">
          Select Page Size:
        </span>
        <select className="select select-sm select-bordered w-16"
          value={page.size}
          onChange={setPageSize}>
          {[10, 20, 50, 100].map((size) => (
            <option key={size} value={size}>
              {size}
            </option>
          ))}
        </select>
      </div>
      <div className="flex gap-2 items-center">
        <button className="btn btn-square"
          aria-label="Previous" title="Previous"
          disabled={page.page === 1}
          onClick={() => setCurrentPage(page.page - 1)}>
          <FaChevronLeft />
        </button>
        <span>
          {page.page} of {totalPages}
        </span>
        <button className="btn btn-square"
          aria-label="Next" title="Next"
          disabled={page.page === totalPages}
          onClick={() => setCurrentPage(page.page + 1)}>
          <FaChevronRight />
        </button>
      </div>
    </footer>
  </div>;
}
