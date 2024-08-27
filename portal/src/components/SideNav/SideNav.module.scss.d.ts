declare namespace SideNavModuleScssNamespace {
  export interface ISideNavModuleScss {
    gettingStarted: string;
    gettingStartedSelected: string;
    link: string;
    selected: string;
    sideNavLinkContainer: string;
    sideNavMainContainer: string;
    text: string;
    userButton: string;
    userButtonContainer: string;
  }
}

declare const SideNavModuleScssModule: SideNavModuleScssNamespace.ISideNavModuleScss;

export = SideNavModuleScssModule;
