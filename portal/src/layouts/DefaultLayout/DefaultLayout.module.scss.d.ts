declare namespace DefaultLayoutModuleScssNamespace {
  export interface IDefaultLayoutModuleScss {
    content: string;
    contentLoading: string;
    contentNoData: string;
    footer: string;
    fullParentHeight: string;
    fullScreenHeight: string;
    header: string;
    layout: string;
    noFooterPadding: string;
    subHeader: string;
  }
}

declare const DefaultLayoutModuleScssModule: DefaultLayoutModuleScssNamespace.IDefaultLayoutModuleScss;

export = DefaultLayoutModuleScssModule;
